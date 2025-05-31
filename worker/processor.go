package worker

import (
	"context"
	"fmt"

	db "github.com/Matltin/Bilitioo-Backend/db/sqlc"
	"github.com/Matltin/Bilitioo-Backend/mail"
	"github.com/hibiken/asynq"
)

const (
	QueueCritical = "critical"
	QueueDefault  = "default"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type RedisTaskProcessor struct {
	server *asynq.Server
	Queries *db.Queries
	mailer mail.EmailSender
}

func NewRedisTaskProcessor(redisOpt asynq.RedisClientOpt, Queries *db.Queries, mailer mail.EmailSender) TaskProcessor {
	server := asynq.NewServer(
		redisOpt,
		asynq.Config{
			Queues: map[string]int{
				QueueCritical: 10,
				QueueDefault:  5,
			},
			ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
				fmt.Println()
			}),
		},
	)

	return &RedisTaskProcessor{
		server: server,
		Queries:  Queries,
		mailer: mailer,
	}
}

func (processor RedisTaskProcessor) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
