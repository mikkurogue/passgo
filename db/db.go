package db

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"

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

	// remove the database store
	os.Remove("./store.db")

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

	sqlStatement := `
  create table foo (id integer not null primary key, name text);
  delete from foo;
  `

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
	stmt, err := tx.Prepare("insert into foo(id, name) values(?,?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	stmt.Exec(1, "test")

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

	rows, err := db.Store.Query("select * from foo")
	if err != nil {
		log.Fatal("HJello world?")
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(id, name)
	}

}
