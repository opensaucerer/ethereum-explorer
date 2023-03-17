# Description: Starts the Indexer and API Server

# Start Indexer
echo "Starting Indexer..."
cd indexer
go build -o ../build/indexer main.go
../build/indexer & # Start Indexer in background

# Start API Server
echo "Starting API Server..."
cd ../api
go build -o ../build/apiserver main.go
../build/apiserver & # Start API Server in background

echo "--------------------------------------------------------------------------------------"
echo "API Service started at http://0.0.0.0:5000"
echo "--------------------------------------------------------------------------------------"