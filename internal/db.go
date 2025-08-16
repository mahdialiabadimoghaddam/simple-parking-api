package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func ConnectToDB() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 dbname=ParkingApp user=postgres password=postgres sslmode=disable connect_timeout=10")
	if err != nil {
		panic(err)
	}

	// Check if the connection actually works
	err = db.Ping()
	if err != nil {

		panic(err)
	}

	return db
}
