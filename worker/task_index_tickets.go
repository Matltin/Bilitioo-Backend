package worker

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
)

const TaskIndexTickets = "task:index_tickets"

type PayloadIndexTickets struct{}

func (distributor *RedisTaskDistributor) DistributeTaskIndexTickets(ctx context.Context, payload *PayloadIndexTickets, opts ...asynq.Option) error {
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskIndexTickets, jsonPayload, opts...)
	info, err := distributor.client.EnqueueContext(ctx, task)
	if err != nil {
		return fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Printf("enqueued task: type=%s, id=%s, queue=%s", info.Type, info.ID, info.Queue)
	return nil
}

func (processor *RedisTaskProcessor) ProcessTaskIndexTickets(ctx context.Context, task *asynq.Task) error {
	tickets, err := processor.Queries.GetAllTickets(ctx)
	if err != nil {
		return fmt.Errorf("failed to get all tickets: %w", err)
	}

	for _, ticket := range tickets {
		doc := map[string]interface{}{
			"origin_city_id":      ticket.OriginProvinceID,
			"destination_city_id": ticket.DestinationProvinceID,
			"vehicle_type":        ticket.VehicleType,
			"departure_time":      ticket.DepartureTime.Format(time.RFC3339),
		}

		data, err := json.Marshal(doc)
		if err != nil {
			log.Printf("Cannot encode ticket %d: %s", ticket.ID, err)
			continue
		}

		_, err = processor.elasticClient.Index(
			"tickets",
			bytes.NewReader(data),
			processor.elasticClient.Index.WithDocumentID(fmt.Sprintf("%d", ticket.ID)),
			processor.elasticClient.Index.WithRefresh("true"),
		)
		if err != nil {
			log.Printf("Error indexing ticket %d: %s", ticket.ID, err)
		}
	}

	log.Printf("indexed %d tickets", len(tickets))
	return nil
}