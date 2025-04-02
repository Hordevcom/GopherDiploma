package storage

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/middleware/logging"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type PGDB struct {
	logger logging.Logger
	DB     *pgxpool.Pool
}

func NewPGDB(logger logging.Logger) *PGDB {
	db, err := pgxpool.New(context.Background(), "postgres://postgres:1@localhost:5432/postgres")

	if err != nil {
		logger.Logger.Errorw("Problem with connecting to db: ", err)
		return nil
	}

	return &PGDB{logger: logger, DB: db}
}
