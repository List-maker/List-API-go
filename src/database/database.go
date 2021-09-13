package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

type Database struct {
	config DatabaseConfig
}

func New(config DatabaseConfig) (*Database, error) {
	db := new(Database)
	db.config = config
	err := db.init()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func (db Database) GetConnection() (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", db.config.Host, db.config.Port, db.config.Username, db.config.Password, db.config.DatabaseName)
	return sql.Open("postgres", connStr)
}
