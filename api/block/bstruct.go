package block

import "github.com/infura/infura-infra-test-perfection-loveday/api/transaction"

type Block struct {
	Number      int64                     `json:"number"`
	Hash        string                    `json:"hash"`
	TxCount     int                       `json:"tx_count"`
	TotalAmount float64                   `json:"total_amount"`
	BlockTime   int64                     `json:"block_time"`
	Txs         []transaction.Transaction `json:"txs"`
}
