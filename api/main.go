package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/infura/infura-infra-test-perfection-loveday/api/constant"
	"github.com/infura/infura-infra-test-perfection-loveday/api/db"
	"github.com/infura/infura-infra-test-perfection-loveday/api/dotenv"
	"github.com/infura/infura-infra-test-perfection-loveday/api/router"
	"github.com/infura/infura-infra-test-perfection-loveday/api/types"
)

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

	m := mux.NewRouter()

	s := &http.Server{
		Addr:     "0.0.0.0:" + constant.Env.Port,
		Handler:  InjectLogging(logger)(m),
		ErrorLog: logger,
	}

	m.HandleFunc("/block", router.HandleLatestBlock).Methods("GET")
	m.HandleFunc("/block/{blockNumber}", router.HandleBlockByNumber).Methods("GET")
	m.HandleFunc("/tx", router.HandleLatestTx).Methods("GET")
	m.HandleFunc("/tx/{txHash}", router.HandleTxByHash).Methods("GET")
	m.HandleFunc("/stats", router.HandleTotalStats).Methods("GET")
	m.HandleFunc("/stats/{blockRange}", router.HandleTotalStatsBetweenBlocks).Methods("GET")
	m.HandleFunc("/index", router.HandleIndexing).Methods("POST")

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

	// connect to the db
	constant.DBClient = db.New(constant.Env.DatabaseURL)
}
