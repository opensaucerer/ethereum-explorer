package transaction

import (
	"context"
	"fmt"
	"log"
	"testing"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/client"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/dotenv"
)

var (
	cl *ethclient.Client
)

func setup() {
	dotenv.Load("../.env")

	cl = client.New(dotenv.Get("INFURA_API_KEY", ""), "rpc")
}

func teardown() {
	cl.Close()
}

func TestTScan(t *testing.T) {

	setup()

	t.Run("Should get transactions in a block", func(t *testing.T) {

		lb, err := cl.BlockByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(lb.Number().String())

		txs, err := GetTransactions(lb)
		if err != nil {
			t.Errorf("Failed to get txs: %v", err)
		}

		if len(txs) != len(lb.Transactions()) {
			t.Errorf("Txs length mismatch: %v != %v", len(txs), len(lb.Transactions()))
		}

	})

	teardown()
}
