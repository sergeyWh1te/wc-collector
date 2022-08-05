package entity

import "github.com/ethereum/go-ethereum/common"

type DepositDTO struct {
	TxHash                string
	BlockNumber		      uint64
	Pubkey                []byte
	WithdrawalCredentials []byte
	Amount                []byte
	Signature             []byte
	Index                 []byte
}

func (d *DepositDTO) GetPubKey() string {
	return common.Bytes2Hex(d.Pubkey[:])
}

func (d *DepositDTO) GetWC() string {
	return common.Bytes2Hex(d.WithdrawalCredentials[:])
}

func (d *DepositDTO) ToValidator(index int64) Validator {
	return Validator{
		TxHash:                d.TxHash,
		Index:                 index,
		PubKey:                d.GetPubKey(),
		WithdrawalCredentials: d.GetWC(),
	}
}
