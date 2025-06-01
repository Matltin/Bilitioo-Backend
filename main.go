package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Matltin/Bilitioo-Backend/api"
	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/mail"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/Matltin/Bilitioo-Backend/worker"
	"github.com/hibiken/asynq"
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

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)
	go runTaskProcessor(config, redisOpt, Queries)

	server := api.NewServer(config, Queries)
	server.Start(":8080")
}

func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store *db.Queries) {
	mail := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAdderss, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mail)
	err := taskProcessor.Start()
	if err != nil {
		fmt.Println("failed to start task processor")
	}
}
