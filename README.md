## Structure

The two microservices:

1. Indexer (blockchain scanner)
2. Public facing REST API

### Indexer

This backend service scans for transactions and captures block information:

- number of transactions
- transaction details (hash, from, to, amount)
- etc.

It is able to achieve this by ping the single endpoint exposed by the indexer. It scans within range if given start-end parameters, or simply subscribes to the latest blocks if no parameters are given.

Examples:

- `/?start=100&end=200` will scan all blocks from 100 to 200 (inclusive)
- `/?start=100` will scan all blocks starting from 100 to the latest, once reached the top it will subscribe to new incoming blocks
- `/` will subscribe to new incoming blocks

### Public API

The REST API service has the following endpoints:

- /block `[GET]`
- /stats `[GET]`
- /tx `[GET]`
- /index `[POST]`

Examples:

- `/block` - returns the latest block and all associated transactions
- `/block/100` - returns the block number 100 and all associated transactions
- `/stats` - returns sum of all amounts and transactions
- `/stats/100:200` - return sum of all amounts and transactions between blocks 100 and 200
- `/tx` - return latest transaction
- `/tx/0x...` - return the transactions with the specified hash
- `index/?auth_token&scan=100:200` - instructs our service to trigger indexer for a fresh scan of blocks between 100 and 200

# Prerequisites

Each microservice comes with a `.env.example` file in their root. This contains the ENV variables that are required for the microservice to run. Copy the contents of this file into a new file called `.env` in the same directory. This file will be ignored by git and will not be committed to the repository. This is where you will store your secrets.

Essentially, you will be need an Infura API key, Postgres DB URL, port numbers, and url to the running indexer node (most likely a localhost url, if they're within the same environment).

# Installation

To get the code up and running, you need to first make a clone to your local hard drive using

```bash
git clone https://github.com/infura/infura-infra-test-perfection-loveday.git explorer
```

Then you need to switch to the right branch

```bash
git checkout mini-blockchain-explorer
```

**After the above steps, it's time to re-visit the requirements section above.**

To run the indexer, run the following command

```bash
cd indexer
go mod tidy
go run main.go
```

The indexer won't automatically start scanning the blockchain. To start scanning, you need to send a POST request to the indexer from the api.

To start the api, open a new terminal instance in the same `explorer` directory and run the following command

```bash
cd api
go mod tidy
go run main.go
```

Now that both the indexer and the api are running, you can send a POST request to the api to start the indexer.

- To scan from block `100` to block `200`

```bash
curl -X POST http://localhost:{your-api-port}/index?auth_token=a-dummy-token&scan=100:200
```

- To scan from block `100` to block `latest` and then keep scanning for new blocks

```bash
curl -X POST http://localhost:{your-api-port}/index?auth_token=a-dummy-token&scan=100:200
```

Note that `auth_token` is required but can include any value.

# Testing

Several parts of the functions used in each microservice have test cases written for them using the Golang standard library. You can run the test cases for each microservice separately.

- Navigate to either the `indexer` directory or `api` directory and run

```bash
go test ./...
```

- To get insights into each test cases, run

```bash
go test -v ./...
```

# Improvements

> Every code can possibly be improved for technology is generative and can be easily describe by a time integral.

1. Use mutex to protect multiple, unnecessary, instantiation of connection.
2. Probably should handle conversion from wei to ether to prevent int bitsize overflow
3. Definitely should add in more test cases
4. Timestamps could be factored into the structs
5. Should add better handler for server shutdown, maybe a quit channel can help (signal.Notify)
6. I surey can clean up the /indexer/main.go file better
7. And of course, some better error handling
