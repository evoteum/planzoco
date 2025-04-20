package main

import (
	"log"
	"planzoco/databases"
	"planzoco/routes"
)

func main() {
	if err := databases.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	r := routes.SetupRoutes()
	r.Run(":8080")
}
