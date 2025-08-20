package worker

import (
	"context"
	"log"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/mail"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
	ProcessTaskIndexTickets(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server      *asynq.Server
	Queries     *db.Queries
	mailer      mail.EmailSender
	elasticClient *elasticsearch.Client
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, Queries *db.Queries, mailer mail.EmailSender, es *elasticsearch.Client) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				log.Printf("Error processing task %s: %v", task.Type(), err)
			}),
		},
	)

	return &RedisTaskProcessor{
		server:      server,
		Queries:     Queries,
		mailer:      mailer,
		elasticClient: es,
	}
}

func (processor RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskIndexTickets, processor.ProcessTaskIndexTickets)
	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	mux.HandleFunc(TaskCleanExpiredReservations, processor.ProcessTaskCleanExpiredReservations)

	return processor.server.Start(mux)
}