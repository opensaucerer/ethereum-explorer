package db

import (
	"context"
	"testing"

	"github.com/infura/infura-infra-test-perfection-loveday/indexer/dotenv"
)

func setup() {
	dotenv.Load("../.env")
}

func TestPG(t *testing.T) {
	setup()

	t.Run("Should init new pg connection", func(t *testing.T) {
		conn := New(dotenv.Get("PG_DB_URL", ""))
		defer conn.Close(context.Background())
	})
}
