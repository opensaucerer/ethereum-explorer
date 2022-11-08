package client

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/ethclient"
)

// create a new eth connection client via infura
func New(key string, conType string) *ethclient.Client {

	log.Println("Connecting to Ethereum node...")

	url := fmt.Sprintf("https://mainnet.infura.io/v3/%s", key)
	if conType == "socket" {
		url = fmt.Sprintf("wss://mainnet.infura.io/ws/v3/%s", key)
	}

	client, err := ethclient.Dial(url)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum node: %v", err)
	}

	log.Println("Connected to Ethereum node")

	return client
}
