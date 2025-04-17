package main

import (
	"log"
    "github.com/evoteum/planzoco/go/planzoco/databases"
    "github.com/evoteum/planzoco/go/planzoco/routes"
)

func main() {
	if err := databases.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := routes.SetupRoutes()
	r.Run(":8080")
}
