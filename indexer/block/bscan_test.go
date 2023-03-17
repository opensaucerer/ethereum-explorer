package block

import (
	"context"
	"log"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/client"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/dotenv"
)

var (
	cl   *ethclient.Client
	step int64 = 10
)

func setup() {
	dotenv.Load("../.env")

	cl = client.New(dotenv.Get("INFURA_API_KEY", ""), "rpc")
}

func teardown() {
	cl.Close()
}

func TestBScan(t *testing.T) {

	setup()

	t.Run("Should get a block by number", func(t *testing.T) {

		lb, err := cl.BlockByNumber(context.Background(), nil)
		if err != nil {
			log.Fatal(err)
		}

		b, err := GetBlock(cl, big.NewInt(lb.Number().Int64()), common.Hash{})
		if err != nil {
			t.Errorf("Failed to get block: %v", err)
		}

		if b.Number != lb.Number().Int64() {
			t.Errorf("Block number mismatch: %v != %v", b.Number, lb.Number().Int64())
		}
	})

	t.Run("Should get blocks by range", func(t *testing.T) {

		lb, err := cl.BlockByNumber(context.Background(), nil)
		if err != nil {
			t.Errorf("Failed to get latest block: %v", err)
		}

		blocks, err := GetBlocks(cl, big.NewInt(lb.Number().Int64()-step), big.NewInt(lb.Number().Int64()))

		if err != nil {
			t.Errorf("Failed to get blocks: %v", err)
		}

		if len(blocks) != int(step+1) {
			t.Errorf("Blocks length mismatch: %v != %v", len(blocks), 10)
		}

		if (blocks[0].Number != lb.Number().Int64()-step) || (blocks[len(blocks)-1].Number != lb.Number().Int64()) {
			t.Errorf("Block number mismatch: %v != %v", blocks[0].Number, lb.Number().Int64()-step)
		}
	})

	teardown()
}
