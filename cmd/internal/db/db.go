package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DataBase struct {
	DB *sql.DB
}

func New() (*DataBase, error) {
	var DB *sql.DB
	var err error

	DB, err = sql.Open("sqlite3", "./app.db")
	return &DataBase{DB: DB}, err
}

func (db *DataBase) Init() {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS collections(
	id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	name TEXT
	);
	`
	_, err := db.DB.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Error Creating the table %q", err)
	}
}
