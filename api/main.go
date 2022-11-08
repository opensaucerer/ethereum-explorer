package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/tech404/backend-challenge/api/dotenv"
	"gitlab.com/tech404/backend-challenge/api/router"
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
		Addr:     ":" + dotenv.Get("PORT", "3000"),
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

	logger.Println("Server started on http://localhost:" + dotenv.Get("PORT", "3000"))

	// create a simple server
	log.Fatal(s.ListenAndServe())
}
