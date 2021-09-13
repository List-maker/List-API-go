package database

import (
	"database/sql"
	"listes_back/src/utils"
)

var commonDb *Database

func SetCommonDb(db *Database) {
	commonDb = db
}

func GetDb() *Database {
	return commonDb
}

func CloseConnection(connection *sql.DB) {
	err := connection.Close()
	if err != nil {
		utils.PrintError(err)
	}
}
