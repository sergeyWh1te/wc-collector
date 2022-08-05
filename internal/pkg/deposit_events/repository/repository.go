package repository

import (
	"context"

	"withdrawal-exchanger/internal/pkg/deposit_events/entity"

	"github.com/jmoiron/sqlx"
)

type repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) Create(ctx context.Context, txHash string, blockNumber int64, pubkey string, index string, withdrawalCredentials string) (*entity.Validator, error) {
	var out entity.Validator

	query := `
insert into node_operators (tx_hash, block_number, pubkey, index, withdrawal_credentials) 
values ($1, $2, $3, $4, $5) 
ON CONFLICT (index) DO NOTHING
returning *;`
	if err := r.db.GetContext(ctx, &out, query, txHash, blockNumber, pubkey, index, withdrawalCredentials); err != nil {
		return nil, err
	}

	return &out, nil
}
