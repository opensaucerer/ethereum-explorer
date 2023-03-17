package types

type Env struct {
	// postgres conn url for db
	DatabaseURL string `env:"PG_DB_URL"`
	// infura api key for ethereum client
	InfuraAPIKey string `env:"INFURA_API_KEY"`
	// port to listen on
	Port string `env:"PORT"`
}
