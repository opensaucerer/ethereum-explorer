package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/block"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/client"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/constant"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/db"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/dotenv"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/subscriber"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/transaction"
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/types"
	"github.com/jackc/pgx/v4"
)

var (
	cl             *ethclient.Client
	pg             *pgx.Conn
	blockChan      = make(chan block.Block, 1000)
	retriesChan    = make(chan block.Block, 1000)
	wg             = &sync.WaitGroup{}
	quitChs        = []chan bool{}
	ticker         = time.NewTicker(1 * time.Second)
	ops            int32
	BG_WORKERS     = 5
	retrials       = map[string]int{}
	retryLimit     = 3
	lock           = &sync.Mutex{}
	mlock          = &sync.Mutex{}
	subscribe      *subscriber.Subscriber
	subControl     ethereum.Subscription
	allowNewBlocks bool
)

func worker(wg *sync.WaitGroup, blockChan chan block.Block, retriesChan chan block.Block, headChan chan *ethtypes.Header, quitCh chan bool) {
	wg.Add(1)
	defer wg.Done()

Loop:
	for {
		select {
		case head := <-headChan:
			if allowNewBlocks {
				log.Printf("New block received: %v\n", head.Number.Int64())
				block, err := block.GetBlock(cl, nil, head.Hash())
				if err != nil {
					log.Printf("Failed to get block %v: %v\n", block.Number, err)
					return
				}
				blockChan <- *block
			}
			// atomic.AddInt32(&ops, 1)
		case err := <-subControl.Err():
			log.Printf("Subscribtion to new blocks failed: %v\n", err)
			break Loop
		case block := <-blockChan:
			// save block to db
			log.Println("Saving block", block.Hash)
			err := block.Save(pg, lock)
			if err != nil {
				log.Printf("Failed to save block %v: %v\n", block.Number, err)
				retriesChan <- block
				atomic.AddInt32(&ops, -1)
				continue
			}
			// get block transactions
			log.Printf("Getting transactions for block %v\n", block.Hash)
			txs, err := transaction.GetTransactions(block.Block)
			if err != nil {
				log.Printf("Some error occurred while getting transactions for block %v: %v\n", block.Number, err)
			}

			for _, tx := range txs {
				log.Printf("Saving transaction %v\n", tx.Hash)
				// save transaction
				err := tx.Save(pg, lock)
				if err != nil {
					log.Printf("Failed to save transaction %v: %v\n", tx.Hash, err)
				}
			}

			log.Printf("Done processing block %v\n", block.Hash)
			atomic.AddInt32(&ops, -1)

		case block := <-retriesChan:
			if retrials[block.Hash] >= retryLimit {
				log.Printf("Block %v has been retried %v times, skipping\n", block.Hash, retryLimit)
				continue
			}
			log.Printf("Retrying block %v\n", block.Hash)
			mlock.Lock()
			retrials[block.Hash]++
			mlock.Unlock()
			blockChan <- block
			atomic.AddInt32(&ops, 1)
			return
		case <-quitCh:
			break Loop
		}
	}
}

func performScan(from, to *big.Int) {
	if to != nil && from.Cmp(to) > 0 {
		return
	}

	// signals to quit all workers
	defer func() {
		for _, quitCh := range quitChs {
			quitCh <- true
		}
		wg.Wait()
	}()

	if to != nil {
		log.Printf("Getting blocks from %v to %v\n", from, to)

		blocks, err := block.GetBlocks(cl, from, to)
		if err != nil {
			log.Printf("Some error occured while getting blocks from %v to %v: %v\n", from, to, err)
		}

		for _, b := range blocks {
			blockChan <- b
			atomic.AddInt32(&ops, 1)
		}
	} else {

		latestBlockNumber, err := cl.BlockNumber(context.Background())
		if err != nil {
			log.Printf("Failed to get latest block number: %v\n", err)
			return
		}

		to = big.NewInt(int64(latestBlockNumber))

		log.Printf("Getting blocks from %v to latest block %v", from, to)

		for i := from; i.Cmp(to) <= 0; i.Add(i, big.NewInt(1)) {
			block, err := block.GetBlock(cl, i, common.Hash{})
			if err != nil {
				log.Printf("Failed to get block %v: %v\n", i, err)
				return
			}
			blockChan <- *block
			atomic.AddInt32(&ops, 1)

			latestBlockNumber, err := cl.BlockNumber(context.Background())
			if err != nil {
				log.Printf("Failed to get latest block number: %v\n", err)
				return
			}
			to = big.NewInt(int64(latestBlockNumber))
		}

		allowNewBlocks = true
		log.Printf("Subscribed to new blocks")

		// set ticker to keep workers infinitely busy
		atomic.AddInt32(&ops, -1)
	}

	for i := 0; i < BG_WORKERS; i++ {
		quitCh := make(chan bool)
		go worker(wg, blockChan, retriesChan, subscribe.HeadChan, quitCh)
		quitChs = append(quitChs, quitCh)
	}

	for tick := range ticker.C {
		if atomic.LoadInt32(&ops) == 0 {
			log.Printf("No more blocks to process, exiting at %v\n", tick)
			break
		}
	}

	log.Println("DONE")
}

func InjectLogging(logger *log.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Println(r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())
			next.ServeHTTP(w, r)
		})
	}
}

func main() {

	logger := log.New(log.Writer(), "http: ", log.LstdFlags)

	logger.Println("Starting server...")

	m := http.NewServeMux()

	s := &http.Server{
		Addr:     "0.0.0.0:" + constant.Env.Port,
		Handler:  InjectLogging(logger)(m),
		ErrorLog: logger,
	}

	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		start := r.URL.Query().Get("start")
		end := r.URL.Query().Get("end")

		if start == "" {
			latestBlockNumber, err := cl.BlockNumber(context.Background())
			if err != nil {
				log.Printf("Failed to get latest block number: %v\n", err)
				return
			}
			start = big.NewInt(int64(latestBlockNumber)).String()
		}

		// convert to int64
		from, err := strconv.ParseInt(start, 10, 64)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"status":  false,
				"message": fmt.Sprintf("failed to parse start index %v", err),
			})
			return
		}

		var to *big.Int
		if end != "" {
			t, err := strconv.ParseInt(end, 10, 64)
			if err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(400)
				json.NewEncoder(w).Encode(map[string]interface{}{
					"status":  false,
					"message": fmt.Sprintf("failed to parse end index %v", err),
				})
				return
			}
			to = big.NewInt(t)
		}

		go performScan(big.NewInt(from), to)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  true,
			"message": "scan started",
		})
	})

	logger.Println("Server started on http://0.0.0.0:" + constant.Env.Port)

	// create a simple server
	log.Fatal(s.ListenAndServe())
}

func init() {
	// load environment variables
	dotenv.Load(".env")

	// verify environment variables or exit program
	err := constant.VerifyEnvironment(types.Env{})
	if err != nil {
		log.Fatal(err)
	}

	// append environment variables to constant.Env
	constant.AppendEnvironment(constant.Env)

	cl = client.New(constant.Env.InfuraAPIKey, "rpc")
	// subscribe to new blocks
	subscribe = subscriber.New(constant.Env.InfuraAPIKey)

	pg = db.New(constant.Env.DatabaseURL)

	err = db.CreateBlocksTable(pg)
	if err != nil {
		log.Printf("Failed to create blocks table: %v\n", err)
	}
	err = db.CreateTransactionsTable(pg)
	if err != nil {
		log.Printf("Failed to create transactions table: %v\n", err)
	}

	// subscribe to new blocks
	subControl, err = subscribe.SubscribeToNewBlocks()
	if err != nil {
		log.Printf("Failed to register to new blocks: %v", err)
	}

}
