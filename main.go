package main

import (
	"database/sql"
	"log"

	"github.com/Matltin/Bilitioo-Backend/api"
	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/util"
	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	DB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	Queries := db.New(DB)

	server := api.NewServer(config, Queries)
	server.Start(":8080")
}
