package transaction

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

func TestTScan(t *testing.T) {

	setup()

	t.Run("Should get a transaction by hash", func(t *testing.T) {
		tx := &Transaction{Hash: "0xfda7aa838107d175945e9cdca4987d67a4ffe78e5c8fce6bec28d2c67b2ba7b8"}
		err := tx.GetByHash(cl, m)
		if err != nil {
			t.Errorf("Failed to get transaction by hash: %v\n", err)
		}
	})

	t.Run("Should get transactions in a block", func(t *testing.T) {

		txs, err := GetTxsByBlockNumber(cl, m, 15833181)
		if err != nil {
			t.Error(err)
		}

		if len(txs) == 0 {
			t.Error("No transactions found")
		}
	})

	t.Run("Should get the latest transaction", func(t *testing.T) {
		tx, err := GetLatestTx(cl, m)
		if err != nil {
			t.Error(err)
		}
		t.Logf("Latest transaction: %v\n", tx)
	})

	teardown()
}
