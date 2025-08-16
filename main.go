package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Matltin/Bilitioo-Backend/api"
	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/mail"
	db_redis "github.com/Matltin/Bilitioo-Backend/redis"
	"github.com/Matltin/Bilitioo-Backend/util"
	"github.com/Matltin/Bilitioo-Backend/worker"
	"github.com/hibiken/asynq"
	_ "github.com/lib/pq"

	_ "github.com/Matltin/Bilitioo-Backend/docs"
	_ "github.com/swaggo/files"       // swagger embed files
	_ "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

//	@title			Bilitio API
//	@version		1.0
//	@description	API for Bilitio ticket booking system.
//	@termsOfService	http://localhost:3000/swagger/index.html

//	@contact.name	API Support
//	@contact.url	http://www.example.com/support
//	@contact.email	support@example.com

//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT

//	@host		localhost:3000
//	@BasePath	/
//	@schemes	http

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
func main() {
	// 1. config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}

	// 2. intialaize db
	DB, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal(err)
	}

	Queries := db.New(DB)

	// 3. intialaize redis
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	// 4. set worker
	// go runTaskProcessor(config, redisOpt, Queries)
	// go runScheduler(redisOpt, taskDistributor)

	redis := db_redis.NewRedisClient(config.RedisAddress)

	// 5. setup server
	server := api.NewServer(config, taskDistributor, Queries, redis)

	server.Start(":3000")
}

// send email task processor
func runTaskProcessor(config util.Config, redisOpt asynq.RedisClientOpt, store *db.Queries) {
	mail := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAdderss, config.EmailSenderPassword)
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store, mail)
	err := taskProcessor.Start()
	if err != nil {
		fmt.Println("failed to start task processor")
	}
}

// check expired reservation processor
func runScheduler(redisOpt asynq.RedisClientOpt, distributor worker.TaskDistributor) {
	scheduler := asynq.NewScheduler(redisOpt, nil)

	_, err := scheduler.Register("* * * * *", worker.NewCleanExpiredReservationsTask())
	if err != nil {
		log.Fatalf("failed to register cron job: %v", err)
	}

	fmt.Println("Scheduler started.")
	if err := scheduler.Run(); err != nil {
		log.Fatalf("failed to run scheduler: %v", err)
	}
}
