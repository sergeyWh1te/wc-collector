package worker

import (
	"context"
	"fmt"
	"sync"

	"withdrawal-exchanger/internal/pkg/deposit_events"
	"withdrawal-exchanger/internal/pkg/deposit_events/entity"
	"withdrawal-exchanger/internal/utils/eth_blocks"

	"github.com/sirupsen/logrus"
)

type worker struct {
	depEvtUc deposit_events.Usecase
	log      *logrus.Logger
}

func New(depEvtUc deposit_events.Usecase, log *logrus.Logger) *worker {
	return &worker{
		depEvtUc: depEvtUc,
		log:      log,
	}
}

const threadsCount = 4

func (w *worker) PrepareChPatrons(from, to, step int64) []<-chan eth_blocks.Patron {
	threads := make([]chan eth_blocks.Patron, threadsCount)
	for i := 0; i < threadsCount; i++ {
		threads[i] = make(chan eth_blocks.Patron)
	}

	go func() {
		defer func() {
			for _, th := range threads {
				close(th)
			}
		}()

		for i, p := range eth_blocks.Patrons(from, to, step) {
			remainderOfDivision := i % 4
			if remainderOfDivision == 0 {
				threads[0] <- p
				continue
			}

			switch float32(remainderOfDivision*100/4) / 100 {
			case 0.25:
				threads[1] <- p
				continue
			case 0.5:
				threads[2] <- p
				continue
			case 0.75:
				threads[3] <- p
				continue
			}
		}
	}()

	out := []<-chan eth_blocks.Patron{
		threads[0],
		threads[1],
		threads[2],
		threads[3],
	}

	return out
}

func (w *worker) Run(ctx context.Context, in ...<-chan eth_blocks.Patron) {
	var wg sync.WaitGroup

	multiplex := func(ctx context.Context, c <-chan eth_blocks.Patron, wn int) {
		defer func() {
			wg.Done()
		}()

		for i := range c {
			w.log.Printf(`Processing blocks %d - %d`, i.From, i.To)

			data, err := w.Handle(ctx, i)
			if err != nil {
				w.log.Errorf(`worker handler error %s`, err)
				continue
			}

			for _, v := range data {
				w.log.Printf(`successfully stored validator index: %d %s`, v.ID, v.PubKey)
			}
		}
	}

	wg.Add(len(in))
	for i, c := range in {
		go multiplex(ctx, c, i)
	}

	wg.Wait()
}

func (w *worker) Handle(ctx context.Context, patron eth_blocks.Patron) ([]entity.Validator, error) {
	eventLog, err := w.depEvtUc.FetchLogs(ctx, patron.From, patron.To)
	if err != nil {
		return nil, fmt.Errorf(`fetchLog error %w`, err)
	}

	depositDTP, err := w.depEvtUc.FilterEvents(eventLog)
	if err != nil {
		return nil, fmt.Errorf(`filter events error %w`, err)
	}

	var out []entity.Validator
	if len(depositDTP) > 0 {
		for result := range w.depEvtUc.ExchangePublicKeyToValidatorIndex(ctx, depositDTP) {
			if result.Error != nil {
				return nil, fmt.Errorf(`pubKey exchange err %s`, result.Error.Error())
			}

			validator, dbErr := w.depEvtUc.Store(ctx,
				result.TxHash,
				int64(result.BlockNumber),
				result.Response.Data.Validator.Pubkey,
				result.Response.Data.Index,
				result.Response.Data.Validator.WithdrawalCredentials,
			)
			if dbErr != nil {
				return nil, fmt.Errorf(`could not save validator into db %w`, dbErr)
			}

			out = append(out, *validator)
		}
	}

	return out, nil
}
