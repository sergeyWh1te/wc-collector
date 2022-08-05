package deposit_events

import (
	"context"

	"withdrawal-exchanger/internal/pkg/deposit_events/entity"

	"github.com/ethereum/go-ethereum/core/types"
)

type Usecase interface {
	FetchLogs(ctx context.Context, fromBlock, toBlock int64) ([]types.Log, error)
	FilterEvents(logs []types.Log) ([]entity.DepositDTO, error)
	ExchangePublicKeyToValidatorIndex(ctx context.Context, deposits []entity.DepositDTO) <-chan entity.Result
	Store(ctx context.Context, txHash string, blockNumber int64, pubkey string, index string, withdrawalCredentials string) (*entity.Validator, error)
}
