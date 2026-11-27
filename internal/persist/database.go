// Package persist contains persistance logic
package persist

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseContext struct {
	Conn *pgxpool.Pool
}

func NewDatabaseContext(ctx context.Context, connString string) *DatabaseContext {
	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("Unable to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database!")
	return &DatabaseContext{Conn: pool}
}
