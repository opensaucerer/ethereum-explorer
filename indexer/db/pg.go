package db

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4"
)

func New(url string) *pgx.Conn {
	log.Printf("Connecting to database...")
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		log.Fatalf("Failed to connect to pg db: %v\n", err)
	}
	log.Printf("Connected to pg db: %v\n", url)
	return conn
}
