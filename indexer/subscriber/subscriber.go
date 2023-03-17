package subscriber

import (
	"context"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/client"
)

// creates an instance of a subcriber
func New(key string) *Subscriber {
	return &Subscriber{
		HeadChan: make(chan *types.Header),
		Client:   client.New(key, "socket"),
	}
}

// registers a subscribtion to new blocks
func (s *Subscriber) SubscribeToNewBlocks() (ethereum.Subscription, error) {
	sub, err := s.Client.SubscribeNewHead(context.Background(), s.HeadChan)
	if err != nil {
		return nil, err
	}
	return sub, nil
}
