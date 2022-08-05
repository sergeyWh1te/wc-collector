package deposit_events

import (
	"context"

	"withdrawal-exchanger/internal/pkg/deposit_events/entity"
)

type Repository interface {
	Create(ctx context.Context, txHash string, blockNumber int64, pubkey string, index string, withdrawalCredentials string) (*entity.Validator, error)
}
