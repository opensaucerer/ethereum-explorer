package types

type Env struct {
	// postgres conn url for db
	DatabaseURL string `env:"PG_DB_URL"`
	// indexer service url
	IndexerURL string `env:"INDEXER_URL"`
	// port to listen on
	Port string `env:"PORT"`
}
