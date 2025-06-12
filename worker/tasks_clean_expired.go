package worker

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
)

const TaskCleanExpiredReservations = "task:clean_expired_reservations"

func NewCleanExpiredReservationsTask() *asynq.Task {
	return asynq.NewTask(TaskCleanExpiredReservations, nil)
}

func (processor *RedisTaskProcessor) ProcessTaskCleanExpiredReservations(ctx context.Context, task *asynq.Task) error {
	err := processor.Queries.MarkExpiredReservations(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete expired reservations: %w", err)
	}

	fmt.Println("✅ expired reservations cleaned successfully")
	return nil
}
