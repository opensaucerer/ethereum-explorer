# Prerequisites

Each microservice comes with an `.env.example` file in their root. This contains the ENV variables that are required for the microservice to run. Copy the contents of this file into a new file called `.env` in the same directory. This file will be ignored by git and will not be committed to the repository. This is where you will store your secrets.

Essentially, you will be need an Infura API key, Postgres DB URL, port numbers, and url to the running indexer node (most likely a localhost url, if they're within the same environment).

# Installation

To get the code up and running, you need to first make a clone to your local hard drive using

```bash
git clone https://gitlab.com/tech404/backend-challenge.git explorer
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
