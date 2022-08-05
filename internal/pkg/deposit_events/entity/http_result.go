package entity

type Result struct {
	BlockNumber uint64
	TxHash      string
	Error       error
	URL         string
	Response    *ValidatorDTO
}
