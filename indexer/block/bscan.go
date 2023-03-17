package block

import (
	"context"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/jackc/pgx/v4"
)

// get a single block by block number
func GetBlock(client *ethclient.Client, blockNumber *big.Int, blockHash common.Hash) (*Block, error) {
	var block *types.Block
	var err error
	if blockNumber != nil {
		block, err = client.BlockByNumber(context.Background(), blockNumber)
		if err != nil {
			return nil, err
		}
	} else {
		block, err = client.BlockByHash(context.Background(), blockHash)
		if err != nil {
			return nil, err
		}
	}
	// calculate total tx amount
	var totalAmount *big.Int = big.NewInt(0)
	for _, tx := range block.Transactions() {
		totalAmount.Add(totalAmount, tx.Value())
	}

	// convert to float64
	amount := new(big.Float).SetInt(totalAmount)
	amountFloat, _ := amount.Float64()

	return &Block{
		Number:      block.Number().Int64(),
		Hash:        block.Hash().Hex(),
		TxCount:     block.Transactions().Len(),
		TotalAmount: amountFloat,
		Block:       block,
		BlockTime:   block.Time(),
	}, nil
}

// insert a block into the db
func (b *Block) Save(db *pgx.Conn, lock *sync.Mutex) error {
	// save to Blocks table
	lock.Lock()
	defer lock.Unlock()
	_, err := db.Exec(context.Background(), `
		INSERT INTO blocks (number, hash, tx_count, total_amount, block_time)
		VALUES ($1, $2, $3, $4, $5)
	`, b.Number, b.Hash, b.TxCount, b.TotalAmount, b.BlockTime)
	if err != nil {
		return err
	}
	return nil
}

// get blocks withing a range of block numbers. if endblocknumber is nil, scan till latest
func GetBlocks(client *ethclient.Client, startBlockNumber, endBlockNumber *big.Int) ([]Block, error) {
	// @Note: this is maybe not the best way to go as it can become memory intensive
	var blocks []Block
	if endBlockNumber.Cmp(startBlockNumber) < 0 {
		return blocks, errors.New("endBlockNumber must be greater than startBlockNumber")
	}
	// scan till latest block if endBlockNumber is nil
	// if endBlockNumber == nil {
	// 	latestBlockNumber, err := client.BlockNumber(context.Background())
	// 	if err != nil {
	// 		return blocks, errors.New("failed to get latest block number: " + err.Error())
	// 	}
	// 	endBlockNumber = big.NewInt(int64(latestBlockNumber))
	// }
	for i := startBlockNumber; i.Cmp(endBlockNumber) <= 0; i.Add(i, big.NewInt(1)) {
		block, err := GetBlock(client, i, common.Hash{})
		if err != nil {
			return nil, err
		}
		blocks = append(blocks, *block)
	}
	return blocks, nil
}
