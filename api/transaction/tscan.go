package transaction

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
)

// retrieve a transaction from the db by tx hash
func (t *Transaction) GetByHash(db *pgx.Conn, lock *sync.Mutex) error {
	lock.Lock()
	defer lock.Unlock()
	err := db.QueryRow(context.Background(), `
		SELECT "block_number", "hash", "from", "to", "amount", "nonce", "block_time"
		FROM transactions
		WHERE "hash" = $1
	`, t.Hash).Scan(&t.BlockNumber, &t.Hash, &t.From, &t.To, &t.Amount, &t.Nonce, &t.BlockTime)
	return err
}

// retrieve transactions from the db by block number
func GetTxsByBlockNumber(db *pgx.Conn, lock *sync.Mutex, blockNumber int64) ([]Transaction, error) {
	lock.Lock()
	defer lock.Unlock()
	rows, err := db.Query(context.Background(), `
		SELECT "block_number", "hash", "from", "to", "amount", "nonce", "block_time"
		FROM transactions
		WHERE "block_number" = $1
	`, blockNumber)

	if err != nil {
		return nil, err
	}

	var txs []Transaction
	for rows.Next() {
		var t Transaction
		err = rows.Scan(&t.BlockNumber, &t.Hash, &t.From, &t.To, &t.Amount, &t.Nonce, &t.BlockTime)
		if err != nil {
			return nil, err
		}
		txs = append(txs, t)
	}
	return txs, nil
}

// retrieve the latest transaction from the db
func GetLatestTx(db *pgx.Conn, lock *sync.Mutex) (Transaction, error) {
	lock.Lock()
	defer lock.Unlock()
	var t Transaction
	err := db.QueryRow(context.Background(), `
		SELECT "block_number", "hash", "from", "to", "amount", "nonce", "block_time"
		FROM transactions
		ORDER BY "block_number" DESC
		LIMIT 1
	`).Scan(&t.BlockNumber, &t.Hash, &t.From, &t.To, &t.Amount, &t.Nonce, &t.BlockTime)
	return t, err
}
