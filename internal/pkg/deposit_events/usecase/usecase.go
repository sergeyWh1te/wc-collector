package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"strings"

	deposit "withdrawal-exchanger/internal/generated/contracts"
	"withdrawal-exchanger/internal/pkg/deposit_events"
	"withdrawal-exchanger/internal/pkg/deposit_events/entity"

	"github.com/avast/retry-go"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type usecase struct {
	client               *ethclient.Client
	repo                 deposit_events.Repository
	addressContract      string
	withdrawalCredential string
	httpClient           *http.Client
	infuraProjectID      string
	infuraPojectPassword string
	infuraHost           string
}

var (
	eventSignature = []byte("DepositEvent(bytes,bytes,bytes,bytes,bytes)")
	depositEvent   = crypto.Keccak256Hash(eventSignature)
)

func New(client *ethclient.Client, repo deposit_events.Repository, addressContract, withdrawalCredential, infuraProjectID, infuraPojectPassword, infuraHost string, httpClient *http.Client) *usecase {
	return &usecase{
		client:               client,
		repo:                 repo,
		addressContract:      addressContract,
		withdrawalCredential: withdrawalCredential,
		infuraProjectID:      infuraProjectID,
		infuraPojectPassword: infuraPojectPassword,
		infuraHost:           infuraHost,
		httpClient:           httpClient,
	}
}

func (u *usecase) FetchLogs(ctx context.Context, fromBlock, toBlock int64) ([]types.Log, error) {
	contractAddress := common.HexToAddress(u.addressContract)
	query := ethereum.FilterQuery{
		FromBlock: big.NewInt(fromBlock),
		ToBlock:   big.NewInt(toBlock),
		Addresses: []common.Address{
			contractAddress,
		},
		Topics: [][]common.Hash{
			{
				depositEvent,
			},
		},
	}

	return u.client.FilterLogs(ctx, query)
}

func (u *usecase) FilterEvents(logs []types.Log) ([]entity.DepositDTO, error) {
	contractAbi, err := abi.JSON(strings.NewReader(deposit.DepositMetaData.ABI))
	if err != nil {
		return nil, err
	}

	var out []entity.DepositDTO
	for _, vLog := range logs {
		if vLog.Topics[0].Hex() != depositEvent.Hex() {
			continue
		}

		depositEventDTO := entity.DepositDTO{
			TxHash:      vLog.TxHash.String(),
			BlockNumber: vLog.BlockNumber,
		}

		if err := contractAbi.UnpackIntoInterface(&depositEventDTO, `DepositEvent`, vLog.Data); err != nil {
			return nil, fmt.Errorf(`could not unmarshal deposit event %w`, err)
		}

		if depositEventDTO.GetWC() == u.withdrawalCredential {
			out = append(out, depositEventDTO)
		}
	}

	return out, nil
}

func (u *usecase) ExchangePublicKeyToValidatorIndex(ctx context.Context, deposits []entity.DepositDTO) <-chan entity.Result {
	results := make(chan entity.Result)
	go func() {
		defer close(results)

		for _, deposit := range deposits {
			url := fmt.Sprintf(`https://%s:%s@%s/eth/v1/beacon/states/finalized/validators/0x%s`, u.infuraProjectID, u.infuraPojectPassword, u.infuraHost, deposit.GetPubKey())

			result := entity.Result{
				BlockNumber: deposit.BlockNumber,
				TxHash:      deposit.TxHash,
				URL:         url,
			}

			result.Response, result.Error = u.fetchValidatorDTO(url)

			select {
			case <-ctx.Done():
				return
			case results <- result:
			}
		}
	}()

	return results
}

func (u *usecase) fetchValidatorDTO(url string) (*entity.ValidatorDTO, error) {
	var out entity.ValidatorDTO

	err := retry.Do(
		func() error {
			resp, err := u.httpClient.Get(url)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			if decodeErr := json.NewDecoder(resp.Body).Decode(&out); decodeErr != nil {
				return fmt.Errorf("could not decode http request: %s; url: %s \n", decodeErr.Error(), url)
			}

			return nil
		},
	)

	return &out, err
}

func (u *usecase) Store(ctx context.Context, txHash string, blockNumber int64, pubkey string, index string, withdrawalCredentials string) (*entity.Validator, error) {
	return u.repo.Create(ctx, txHash, blockNumber, pubkey, index, withdrawalCredentials)
}
