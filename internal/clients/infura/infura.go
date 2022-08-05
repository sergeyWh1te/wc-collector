package infura

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
)

func NewHttp(ctx context.Context, infuraKEY string) (*ethclient.Client, error) {
	return ethclient.DialContext(ctx, fmt.Sprintf(`https://mainnet.infura.io/v3/%s`, infuraKEY))
}
