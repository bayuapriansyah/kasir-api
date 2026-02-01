package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) {
	//buka database
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	//test koneksi
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	//set koneksi pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(0)

	log.Println("Berhasil koneksi ke database")
	return db, nil
}
