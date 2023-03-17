FROM golang:1.19.0-alpine3.15

RUN apk add build-base

RUN mkdir -p /go/src/github.com/infura/infura-infra-test-perfection-loveday/indexer

WORKDIR /go/src/github.com/infura/infura-infra-test-perfection-loveday/indexer

COPY go.mod .
COPY go.sum .

ENV PORT $PORT
ENV PG_DB_URL $PG_DB_URL
ENV INFURA_API_KEY $INFURA_API_KEY

RUN go mod download

COPY . .

RUN go build -o indexer

RUN chmod +x indexer

EXPOSE $PORT

ENTRYPOINT ./indexer