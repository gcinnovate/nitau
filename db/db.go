package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //import postgres

	"github.com/gcinnovate/nitau/config"
)

var db *sqlx.DB

func init() {
	psqlInfo := fmt.Sprintf("%s", config.NitauConf.Database.URI)

	var err error
	db, err = ConnectDB(psqlInfo)
	if err != nil {
		log.Fatal(err)
	}
}

// ConnectDB ...
func ConnectDB(dataSourceName string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", dataSourceName)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return db, nil
}

//GetDB ...
func GetDB() *sqlx.DB {
	return db
}
