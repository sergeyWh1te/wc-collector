package eth_blocks

type Patron struct {
	From int64
	To   int64
}

func Patrons(fromBlock, toBlock, step int64) []Patron {
	var patrons []Patron
	var carry int64
	for carry = fromBlock; carry < toBlock - step; carry += step {
		patrons = append(patrons, Patron{
			From: carry,
			To:   carry + step,
		})
	}

	if toBlock > carry {
		patrons = append(patrons, Patron{
			From: carry,
			To:   toBlock,
		})
	}

	return patrons
}