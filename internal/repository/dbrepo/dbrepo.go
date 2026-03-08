package dbrepo

import (
	"database/sql"

	"github.com/bangn/bookings/internal/config"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresRepo(a *config.AppConfig, conn *sql.DB) *PostgresDBRepo {
	return &PostgresDBRepo{
		App: a,
		DB: conn,
	}
}