package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var PgxPool *pgxpool.Pool

func InitPostgres() {
	postgresURL := os.Getenv("POSTGRES_URL")
	if postgresURL == "" {
		log.Fatal("POSTGRES_URL tidak diatur di file .env")
	}

	var err error
	PgxPool, err = pgxpool.New(context.Background(), postgresURL)
	if err != nil {
		log.Fatalf("Tidak dapat membuat connection pool ke PostgreSQL: %v\n", err)
	}

	if err := PgxPool.Ping(context.Background()); err != nil {
		log.Fatalf("Tidak dapat terhubung ke PostgreSQL: %v\n", err)
	}

	log.Printf("Terhubung ke PostgreSQL!")
}

func ClosePostgres() {
	if PgxPool != nil {
		PgxPool.Close()
		fmt.Println("Koneksi PostgreSQL ditutup.")
	}
}
