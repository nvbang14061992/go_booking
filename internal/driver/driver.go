package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// DB is the database connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL creates a connection pool for the database and returns the DB struct
func ConnectSQL(dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open connections, idle connections, and connection lifetime
	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifetime)

	// Assign the database connection to the global variable
	dbConn.SQL = db
	err = testDB(dbConn.SQL)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// testDB pings the database to make sure we have a connection pool established
func testDB(d *sql.DB) (error) {
	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}


// NewDatabase creates a connection pool for the database and returns the DB struct
func NewDatabase(dsn string) (*DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}
	return &DB{SQL: db}, nil
}