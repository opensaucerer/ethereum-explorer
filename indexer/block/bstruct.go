package block

import "github.com/ethereum/go-ethereum/core/types"

// CREATE TABLE blocks(
//     number INT PRIMARY KEY,
//     hash CHAR(32) NOT NULL,
//     tx_count INT NOT NULL
// );

type Block struct {
	Number      int64   `json:"number"`
	Hash        string  `json:"hash"`
	TxCount     int     `json:"tx_count"`
	TotalAmount float64 `json:"total_amount"`
	BlockTime   uint64  `json:"block_time"`
	Block       *types.Block
}
