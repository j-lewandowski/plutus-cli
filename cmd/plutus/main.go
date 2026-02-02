package main

import (
	"fmt"
	"plutus-cli/internal/db"
)

func main() {
	repository, err := db.InitDb()

	if err != nil {
		fmt.Println("Error initializing database:", err)
		return
	}
	defer repository.Close()

	Execute(repository)
}
