package db

import (
	"errors"
	"log"
	"os"
)

func Bootstrap(db Database) {

	if db.Store == nil {
		log.Fatal("No valid connection to the data store")
	}

	RemoveStore()

	connErr := db.CreateInitialConnection()
	if connErr != nil {
		log.Fatal(connErr.Error())
	}

	err := db.CreateStoreTable()
	if err != nil {
		log.Fatal(err.Error())
	}
}

func RemoveStore() {
	os.Remove("./store.db")
}

func CheckIfStoreExists() bool {
	if _, err := os.Stat("./store.db"); errors.Is(err, os.ErrNotExist) {
		return false
	}

	return true
}
