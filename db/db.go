package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Connected bool
	Store     *sql.DB
}

type Service struct {
	Username string
	Password string
	Service  string // The app/site/whatever you want to store
}

func (db *Database) CreateInitialConnection() error {

	store, err := sql.Open("sqlite3", "./store.db")
	if err != nil {
		return err
	}

	db.Store = store
	db.Connected = true

	fmt.Println("Connected?")
	return nil
}

func (db *Database) CloseConnection() error {
	if db.Connected == false {
		return errors.New("Connection never established")
	}

	defer db.Store.Close()
	return nil
}

func (db Database) CreateStoreTable() error {

	if db.Store == nil {
		log.Fatal("No active connection to data store")
	}

	sqlStatement := CREATE_TABLE

	_, err := db.Store.Exec(sqlStatement)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	return nil
}

func (db *Database) InsertService(service Service) error {

	if db.Store == nil {
		log.Fatal("No active connection to data store")
	}

	// start transaction
	tx, err := db.Store.Begin()
	if err != nil {
		return err
	}

	// prepare query
	stmt, err := tx.Prepare(INSERT_SERVICE)
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(service.Username, service.Password, service.Service)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GetAllServices() {

	if db.Store == nil {
		log.Fatal("No active connection to the data store")
	}

	rows, err := db.Store.Query("select * from services")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username, password, service string

		err = rows.Scan(&id, &username, &password, &service)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, username, password, service)
	}

}
