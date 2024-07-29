package initdb

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDb(dataSourceUrl string) error {
	var err error
	DB, err = sql.Open("postgres", dataSourceUrl)
	if err != nil {
		return err
	}
	return DB.Ping()
}
