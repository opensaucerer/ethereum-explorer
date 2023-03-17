package constant

import (
	"github.com/infura/infura-infra-test-perfection-loveday/indexer/types"
)

const (
	envTagName      = "env" // environment variable tag name
	ShutdownTimeout = 5     // seconds
)

var (
	Env = new(types.Env) // global environment variable
)
