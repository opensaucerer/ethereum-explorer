package block

import (
	"context"
	"sync"
	"testing"

	"github.com/infura/infura-infra-test-perfection-loveday/api/db"
	"github.com/infura/infura-infra-test-perfection-loveday/api/dotenv"
	"github.com/jackc/pgx/v4"
)

var (
	cl *pgx.Conn
	m  = &sync.Mutex{}
)

func setup() {
	dotenv.Load("../.env")

	cl = db.New(dotenv.Get("PG_DB_URL", ""))
}

func teardown() {
	cl.Close(context.Background())
}

func TestBScan(t *testing.T) {

	setup()

	t.Run("Should get a block by number", func(t *testing.T) {
		b := &Block{Number: 15833181}
		err := b.GetByNumber(cl, m)
		if err != nil {
			t.Errorf("Failed to get block by number: %v\n", err)
		}
	})

	t.Run("Should get a block by hash", func(t *testing.T) {
		b := &Block{Hash: "0x774867b353426e8d3e5054f6f06a1605ab9731b0b61761232b610f46d3f01633"}
		err := b.GetByHash(cl, m)
		if err != nil {
			t.Errorf("Failed to get block by hash: %v\n", err)
		}
	})

	t.Run("Should get total stats", func(t *testing.T) {
		txCount, totalAmount, err := GetTotalStats(cl, m)
		if err != nil {
			t.Errorf("Failed to get total stats: %v\n", err)
		}
		t.Logf("Total transactions: %v\n", txCount)
		t.Logf("Total amount: %v\n", totalAmount)
	})

	t.Run("Should get total stats between two blocks", func(t *testing.T) {
		txCount, totalAmount, err := GetStatsByBlockNumber(cl, m, 15833181, 15833183)
		if err != nil {
			t.Errorf("Failed to get total stats between two blocks: %v\n", err)
		}
		t.Logf("Total transactions: %v\n", txCount)
		t.Logf("Total amount: %v\n", totalAmount)
	})

	teardown()
}
