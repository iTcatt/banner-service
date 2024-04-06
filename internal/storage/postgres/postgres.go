package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
)

type Storage struct {
	conn *pgx.Conn
}

func NewStorage(dbPath string) (*Storage, error) {
	conn, err := pgx.Connect(context.Background(), dbPath)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, err
	}
	log.Println("successful database connection")

	return &Storage{conn: conn}, nil
}

func (s *Storage) Close(_ context.Context) error {
	return nil
}
