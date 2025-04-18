package main

import (
	"log"
    "github.com/evoteum/planzoco/go/planzoco/databases"
    "github.com/evoteum/planzoco/go/planzoco/routes"
)

func main() {
	ctx := context.Background()
	if err := databases.InitDB(ctx); err != nil {
		log.Fatal(err)
	}

	r := routes.SetupRoutes()
	r.Run(":8080")
}
