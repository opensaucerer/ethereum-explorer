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
	"github.com/jackc/pgx/v4"
	"gitlab.com/tech404/backend-challenge/api/block"
	"gitlab.com/tech404/backend-challenge/api/db"
	"gitlab.com/tech404/backend-challenge/api/dotenv"
	"gitlab.com/tech404/backend-challenge/api/transaction"
)

var (
	cl   *pgx.Conn
	lock = &sync.Mutex{}
)

func init() {
	dotenv.Load("")
	// connect to the db
	cl = db.New(dotenv.Get("PG_DB_URL", ""))
}

func HandleLatestBlock(w http.ResponseWriter, r *http.Request) {
	// get the latest block from the db
	block, err := block.GetLatestBlock(cl, lock)
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
	txs, err := transaction.GetTxsByBlockNumber(cl, lock, block.Number)
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
	err = block.GetByNumber(cl, lock)
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
	txs, err := transaction.GetTxsByBlockNumber(cl, lock, block.Number)
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
	tx, err := transaction.GetLatestTx(cl, lock)
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
	err := tx.GetByHash(cl, lock)
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
	txCount, totalAmount, err := block.GetTotalStats(cl, lock)
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
	txCount, totalAmount, err := block.GetStatsByBlockNumber(cl, lock, startBlock, endBlock)
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

	scan := r.URL.Query().Get("scan")
	if scan == "" {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  false,
			"message": "scan is required. e.g scan=100:200",
		})
		return
	}

	if strings.Contains(scan, ":") {
		start := strings.Split(scan, ":")[0]
		end := strings.Split(scan, ":")[1]
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
		return

	} else {
		start, err := strconv.ParseInt(scan, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  false,
				"message": "invalid block number",
			})
			return
		}

		// make an http request to the indexer
		resp, err := http.Get(fmt.Sprintf(dotenv.Get("INDEXER_URL", "")+"?start=%v", start))
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
		return
	}

}
