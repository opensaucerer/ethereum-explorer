package constant

import (
	"github.com/infura/infura-infra-test-perfection-loveday/api/types"
	"github.com/jackc/pgx/v4"
)

const (
	envTagName      = "env" // environment variable tag name
	ShutdownTimeout = 5     // seconds
)

var (
	Env      = new(types.Env) // global environment variable
	DBClient *pgx.Conn        // global database client
)
