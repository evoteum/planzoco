package main

import (
	"log"
	"github.com/evoteum/planzoco/databases"
	"github.com/evoteum/planzoco/routes"
)

func main() {
	if err := databases.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := routes.SetupRoutes()
	r.Run(":8080")
}
