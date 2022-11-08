package transaction

import (
	"context"
	"log"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/jackc/pgx/v4"
)

// get all transactions in a block
func GetTransactions(block *types.Block) ([]Transaction, error) {
	var txs []Transaction
	var err error
	for _, tx := range block.Transactions() {
		var signer common.Address
		signer, err = types.LatestSignerForChainID(tx.ChainId()).Sender(tx)
		if err != nil {
			log.Printf("Failed to get sender address for tx %s: %v\n", tx.Hash().Hex(), err)
			continue
		}
		to := ""
		if tx.To() != nil {
			to = tx.To().Hex()
		}
		from := ""
		if signer != (common.Address{}) {
			from = signer.Hex()
		}

		txs = append(txs, Transaction{
			BlockNumber: block.Number().Int64(),
			Hash:        tx.Hash().Hex(),
			From:        from,
			To:          to,
			Amount:      float64(tx.Value().Int64()),
			Nonce:       tx.Nonce(),
			BlockTime:   block.Time(),
		})
	}
	return txs, err
}

// save transaction to db
func (t *Transaction) Save(db *pgx.Conn, lock *sync.Mutex) error {
	// save to Transactions table
	lock.Lock()
	defer lock.Unlock()
	_, err := db.Exec(context.Background(), `
		INSERT INTO transactions ("block_number", "hash", "from", "to", "amount", "nonce", "block_time")
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, t.BlockNumber, t.Hash, t.From, t.To, t.Amount, t.Nonce, t.BlockTime)
	if err != nil {
		return err
	}
	return nil
}
