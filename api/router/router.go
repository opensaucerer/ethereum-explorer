package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/gorilla/mux"
	"github.com/infura/infura-infra-test-perfection-loveday/api/block"
	"github.com/infura/infura-infra-test-perfection-loveday/api/constant"
	"github.com/infura/infura-infra-test-perfection-loveday/api/dotenv"
	"github.com/infura/infura-infra-test-perfection-loveday/api/transaction"
)

var (
	lock = &sync.Mutex{}
)

func HandleLatestBlock(w http.ResponseWriter, r *http.Request) {
	// get the latest block from the db
	block, err := block.GetLatestBlock(constant.DBClient, lock)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to get latest block " + err.Error(),
		})
		return
	}

	// get the transactions in the block
	txs, err := transaction.GetTxsByBlockNumber(constant.DBClient, lock, block.Number)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to get latest block " + err.Error(),
		})
		return
	}

	block.Txs = txs

	// return the block as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"data":   block,
	})
}

func HandleBlockByNumber(w http.ResponseWriter, r *http.Request) {
	// get the block number from the url
	vars := mux.Vars(r)
	blockNumber, err := strconv.ParseInt(vars["blockNumber"], 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "invalid block number",
		})
		return
	}

	// get the block from the db
	block := &block.Block{Number: blockNumber}
	err = block.GetByNumber(constant.DBClient, lock)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to get block by number " + err.Error(),
		})
		return
	}

	// get the transactions in the block
	txs, err := transaction.GetTxsByBlockNumber(constant.DBClient, lock, block.Number)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to get block by number " + err.Error(),
		})
		return
	}
	block.Txs = txs
	// return the block as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"data":   block,
	})
}

func HandleLatestTx(w http.ResponseWriter, r *http.Request) {
	// get the latest transaction from the db
	tx, err := transaction.GetLatestTx(constant.DBClient, lock)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to get latest transaction " + err.Error(),
		})
		return
	}

	// return the transaction as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"data":   tx,
	})
}

func HandleTxByHash(w http.ResponseWriter, r *http.Request) {
	// get the transaction hash from the url
	vars := mux.Vars(r)
	txHash := vars["txHash"]

	// get the transaction from the db
	tx := &transaction.Transaction{Hash: txHash}
	err := tx.GetByHash(constant.DBClient, lock)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to get transaction by hash " + err.Error(),
		})
		return
	}

	// return the transaction as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"data":   tx,
	})
}

func HandleTotalStats(w http.ResponseWriter, r *http.Request) {
	// get the total stats from the db
	txCount, totalAmount, err := block.GetTotalStats(constant.DBClient, lock)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to get total stats " + err.Error(),
		})
		return
	}

	// return the stats as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"data": map[string]interface{}{
			"tx_count":     txCount,
			"total_amount": totalAmount,
		},
	})
}

func HandleTotalStatsBetweenBlocks(w http.ResponseWriter, r *http.Request) {
	// get the block numbers from the url
	vars := mux.Vars(r)
	blockRange := vars["blockRange"]
	start := strings.Split(blockRange, ":")[0]
	end := strings.Split(blockRange, ":")[1]
	startBlock, err := strconv.ParseInt(start, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "invalid start block number",
		})
		return
	}
	endBlock, err := strconv.ParseInt(end, 10, 64)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "invalid end block number",
		})
		return
	}

	if startBlock > endBlock {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "start block number cannot be greater than end block number",
		})
		return
	}

	// get the total stats from the db
	txCount, totalAmount, err := block.GetStatsByBlockNumber(constant.DBClient, lock, startBlock, endBlock)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to get total stats between blocks " + err.Error(),
		})
		return
	}

	// return the stats as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": true,
		"data": map[string]interface{}{
			"tx_count":     txCount,
			"total_amount": totalAmount,
		},
	})
}

func HandleIndexing(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("auth_token")
	if token == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "auth token is required",
		})
		return
	}

	var start string
	var end string
	var err error

	scan := r.URL.Query().Get("scan")
	if scan != "" {

		if strings.Contains(scan, ":") {
			start = strings.Split(scan, ":")[0]
			end = strings.Split(scan, ":")[1]
			startBlock, err := strconv.ParseInt(start, 10, 64)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status":  false,
					"message": "invalid start block number",
				})
				return
			}
			endBlock, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status":  false,
					"message": "invalid end block number",
				})
				return
			}

			if startBlock > endBlock {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status":  false,
					"message": "start block number cannot be greater than end block number",
				})
				return
			}

		} else {
			start = scan
			_, err := strconv.ParseInt(start, 10, 64)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status":  false,
					"message": "invalid block number",
				})
				return
			}

		}
	}

	// make an http request to the indexer
	resp, err := http.Get(fmt.Sprintf(dotenv.Get("INDEXER_URL", "")+"?start=%s&end=%s", start, end))
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to make request to indexer " + err.Error(),
		})
		return
	}

	// read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "failed to read response body from indexer " + err.Error(),
		})
		return
	}

	// return the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)

}
