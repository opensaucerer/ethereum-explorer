# Description: Starts the Indexer and API Server

# Start Indexer
echo "Starting Indexer..."
cd indexer
go build -o indexer main.go
./indexer & # Start Indexer in background

# Start API Server
echo "Starting API Server..."
cd ../api
go build -o apiserver main.go
./apiserver & # Start API Server in background