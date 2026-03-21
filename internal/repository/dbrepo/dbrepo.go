package dbrepo

import (
	"database/sql"

	"github.com/bangn/bookings/internal/config"
	"github.com/bangn/bookings/internal/repository"
)

type PostgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

type testDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresRepo(a *config.AppConfig, conn *sql.DB) repository.DatabaseRepo {
	return &PostgresDBRepo{
		App: a,
		DB: conn,
	}
}

func NewTestingPostgresRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo {
		App: a,
	}
}