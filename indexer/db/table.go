package db

import (
	"context"

	"github.com/jackc/pgx/v4"
)

func CreateBlocksTable(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS blocks (
			number INT PRIMARY KEY,
			hash VARCHAR NOT NULL,
			tx_count INT NOT NULL,
			total_amount NUMERIC NOT NULL,
			block_time NUMERIC NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}

func CreateTransactionsTable(conn *pgx.Conn) error {
	_, err := conn.Exec(context.Background(), `
		CREATE TABLE IF NOT EXISTS transactions (
			"hash" VARCHAR PRIMARY KEY,
			"block_number" INT NOT NULL,
			"from" VARCHAR NOT NULL,
			"to" VARCHAR NOT NULL,
			"amount" NUMERIC NOT NULL,
			"nonce" INT NOT NULL,
			"block_time" NUMERIC NOT NULL
		);
	`)
	if err != nil {
		return err
	}
	return nil
}
