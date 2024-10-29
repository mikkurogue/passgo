package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// entry point
func main() {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("No env file found... exiting...")
		os.Exit(0)
	}
}
