FROM golang:1.19.0-alpine3.15

RUN mkdir -p /go/src/github.com/infura/infura-infra-test-perfection-loveday/api

WORKDIR /go/src/github.com/infura/infura-infra-test-perfection-loveday/api

COPY go.mod .
COPY go.sum .

ENV PORT $PORT
ENV PG_DB_URL $PG_DB_URL
ENV INDEXER_URL $INDEXER_URL

RUN go mod download

COPY . .

RUN go build -o api

RUN chmod +x api

EXPOSE $PORT

ENTRYPOINT ./api