package client

import (
	"testing"

	"github.com/infura/infura-infra-test-perfection-loveday/indexer/dotenv"
)

func setup() {
	dotenv.Load("../.env")
}

func TestEth(t *testing.T) {

	setup()

	t.Run("Should init new ethereum connection", func(t *testing.T) {
		client := New(dotenv.Get("INFURA_API_KEY", ""), "rpc")
		defer client.Close()
	})
}
