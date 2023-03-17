package block

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4"
)

// retrieve a block from the db by block hash
func (b *Block) GetByHash(db *pgx.Conn, lock *sync.Mutex) error {
	lock.Lock()
	rows, err := db.Query(context.Background(), `
		SELECT "number", "hash", "tx_count", "total_amount", "block_time"
		FROM blocks
		WHERE "hash" = $1
	`, b.Hash)
	defer lock.Unlock()

	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&b.Number, &b.Hash, &b.TxCount, &b.TotalAmount, &b.BlockTime)
		if err != nil {
			return err
		}
	}
	return nil
}

// retrieve a block from the db by block number
func (b *Block) GetByNumber(db *pgx.Conn, lock *sync.Mutex) error {
	lock.Lock()
	defer lock.Unlock()
	rows, err := db.Query(context.Background(), `
		SELECT "number", "hash", "tx_count", "total_amount", "block_time"
		FROM blocks
		WHERE "number" = $1
	`, b.Number)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&b.Number, &b.Hash, &b.TxCount, &b.TotalAmount, &b.BlockTime)
		if err != nil {
			return err
		}
	}
	return nil
}

// sum of all amount and transactions in db
func GetTotalStats(db *pgx.Conn, lock *sync.Mutex) (int64, float64, error) {
	lock.Lock()
	defer lock.Unlock()
	rows, err := db.Query(context.Background(), `
		SELECT SUM("tx_count"), SUM("total_amount")
		FROM blocks
	`)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var txCount int64
	var totalAmount float64
	for rows.Next() {
		err = rows.Scan(&txCount, &totalAmount)
		if err != nil {
			return 0, 0, err
		}
	}
	return txCount, totalAmount, nil
}

// sum of all amount and transactions between two block numbers
func GetStatsByBlockNumber(db *pgx.Conn, lock *sync.Mutex, startBlockNumber, endBlockNumber int64) (int64, float64, error) {
	lock.Lock()
	defer lock.Unlock()
	rows, err := db.Query(context.Background(), `
		SELECT SUM("tx_count"), SUM("total_amount")
		FROM blocks
		WHERE "number" BETWEEN $1 AND $2
	`, startBlockNumber, endBlockNumber)
	if err != nil {
		return 0, 0, err
	}
	defer rows.Close()

	var txCount int64
	var totalAmount float64
	for rows.Next() {
		err = rows.Scan(&txCount, &totalAmount)
		if err != nil {
			return 0, 0, err
		}
	}
	return txCount, totalAmount, nil
}

// retrieve latest block from the db
func GetLatestBlock(db *pgx.Conn, lock *sync.Mutex) (Block, error) {
	lock.Lock()
	defer lock.Unlock()
	rows, err := db.Query(context.Background(), `
		SELECT "number", "hash", "tx_count", "total_amount", "block_time"
		FROM blocks
		ORDER BY "number" DESC
		LIMIT 1
	`)
	if err != nil {
		return Block{}, err
	}
	defer rows.Close()

	var b Block
	for rows.Next() {
		err = rows.Scan(&b.Number, &b.Hash, &b.TxCount, &b.TotalAmount, &b.BlockTime)
		if err != nil {
			return Block{}, err
		}
	}
	return b, nil
}
