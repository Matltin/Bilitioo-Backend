package main

import (
	"database/sql"
	"log"

	"github.com/Matltin/Bilitioo-Backend/api"
	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	_ "github.com/lib/pq"
)

func main() {
	DB, err := sql.Open("postgres", "postgresql://root:secret@localhost:5432/bilitioo?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	Queries := db.New(DB)

	server := api.NewServer(Queries)
	server.Start(":8080")
}
