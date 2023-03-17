package subscriber

import (
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Subscriber struct {
	HeadChan chan *types.Header
	Client   *ethclient.Client
}
