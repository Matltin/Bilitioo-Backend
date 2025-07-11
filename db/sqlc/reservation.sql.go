// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: reservation.sql

package db

import (
	"context"
	"time"
)

const cancelReservation = `-- name: CancelReservation :exec
UPDATE "ticket"
SET status = 'NOT_RESERVED'
WHERE id = $1 AND status = 'RESERVED'
`

func (q *Queries) CancelReservation(ctx context.Context, id int64) error {
	_, err := q.db.ExecContext(ctx, cancelReservation, id)
	return err
}

const createReservation = `-- name: CreateReservation :one
INSERT INTO "reservation" (
    "user_id",
    "ticket_id",
    "payment_id"
) VALUES (
    $1, $2, $3
) RETURNING 
    id, 
    user_id, 
    ticket_id, 
    payment_id, 
    status, 
    EXTRACT(EPOCH FROM duration_time)::bigint as duration_time_seconds,
    created_at
`

type CreateReservationParams struct {
	UserID    int64 `json:"user_id"`
	TicketID  int64 `json:"ticket_id"`
	PaymentID int64 `json:"payment_id"`
}

type CreateReservationRow struct {
	ID                  int64        `json:"id"`
	UserID              int64        `json:"user_id"`
	TicketID            int64        `json:"ticket_id"`
	PaymentID           int64        `json:"payment_id"`
	Status              TicketStatus `json:"status"`
	DurationTimeSeconds int64        `json:"duration_time_seconds"`
	CreatedAt           time.Time    `json:"created_at"`
}

func (q *Queries) CreateReservation(ctx context.Context, arg CreateReservationParams) (CreateReservationRow, error) {
	row := q.db.QueryRowContext(ctx, createReservation, arg.UserID, arg.TicketID, arg.PaymentID)
	var i CreateReservationRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TicketID,
		&i.PaymentID,
		&i.Status,
		&i.DurationTimeSeconds,
		&i.CreatedAt,
	)
	return i, err
}

const getAllUserReservation = `-- name: GetAllUserReservation :many
SELECT 
    re.id,
    t.id,
    oc.province,
    dc.province
FROM "reservation" re 
INNER JOIN "ticket" t ON re.ticket_id = t.id
INNER JOIN "route" ro ON t.route_id = ro.id
INNER JOIN "city" oc ON oc.id = ro.origin_city_id
INNER JOIN "city" dc ON dc.id = ro.destination_city_id
WHERE re.status != 'RESERVED' AND re.user_id = $1
`

type GetAllUserReservationRow struct {
	ID         int64  `json:"id"`
	ID_2       int64  `json:"id_2"`
	Province   string `json:"province"`
	Province_2 string `json:"province_2"`
}

func (q *Queries) GetAllUserReservation(ctx context.Context, userID int64) ([]GetAllUserReservationRow, error) {
	rows, err := q.db.QueryContext(ctx, getAllUserReservation, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetAllUserReservationRow{}
	for rows.Next() {
		var i GetAllUserReservationRow
		if err := rows.Scan(
			&i.ID,
			&i.ID_2,
			&i.Province,
			&i.Province_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getCompletedUserReservation = `-- name: GetCompletedUserReservation :many
SELECT 
    re.id AS "reservation_id",
    t.id AS "ticket_id",
    oc.province,
    dc.province
FROM "reservation" re 
INNER JOIN "ticket" t ON re.ticket_id = t.id
INNER JOIN "route" ro ON t.route_id = ro.id
INNER JOIN "city" oc ON oc.id = ro.origin_city_id
INNER JOIN "city" dc ON dc.id = ro.destination_city_id
WHERE re.status = 'RESERVED' AND re.user_id = $1
`

type GetCompletedUserReservationRow struct {
	ReservationID int64  `json:"reservation_id"`
	TicketID      int64  `json:"ticket_id"`
	Province      string `json:"province"`
	Province_2    string `json:"province_2"`
}

func (q *Queries) GetCompletedUserReservation(ctx context.Context, userID int64) ([]GetCompletedUserReservationRow, error) {
	rows, err := q.db.QueryContext(ctx, getCompletedUserReservation, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []GetCompletedUserReservationRow{}
	for rows.Next() {
		var i GetCompletedUserReservationRow
		if err := rows.Scan(
			&i.ReservationID,
			&i.TicketID,
			&i.Province,
			&i.Province_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getIDReservation = `-- name: GetIDReservation :many
SELECT r.id FROM "reservation" r
INNER JOIN "payment" p ON p.id = r.payment_id
WHERE p.id = $1
`

func (q *Queries) GetIDReservation(ctx context.Context, id int64) ([]int64, error) {
	rows, err := q.db.QueryContext(ctx, getIDReservation, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []int64{}
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getReservationDetails = `-- name: GetReservationDetails :one
SELECT
    r.id,
    r.payment_id,
    t.amount, 
    r.user_id,
    t.departure_time,
    t.status
FROM "ticket" t
INNER JOIN "reservation" r ON r.ticket_id = t.id 
WHERE t.id = $1
`

type GetReservationDetailsRow struct {
	ID            int64                        `json:"id"`
	PaymentID     int64                        `json:"payment_id"`
	Amount        int64                        `json:"amount"`
	UserID        int64                        `json:"user_id"`
	DepartureTime time.Time                    `json:"departure_time"`
	Status        CheckReservationTicketStatus `json:"status"`
}

func (q *Queries) GetReservationDetails(ctx context.Context, id int64) (GetReservationDetailsRow, error) {
	row := q.db.QueryRowContext(ctx, getReservationDetails, id)
	var i GetReservationDetailsRow
	err := row.Scan(
		&i.ID,
		&i.PaymentID,
		&i.Amount,
		&i.UserID,
		&i.DepartureTime,
		&i.Status,
	)
	return i, err
}

const getReservationStatus = `-- name: GetReservationStatus :one
SELECT status FROM "reservation" 
WHERE id = $1
`

func (q *Queries) GetReservationStatus(ctx context.Context, id int64) (TicketStatus, error) {
	row := q.db.QueryRowContext(ctx, getReservationStatus, id)
	var status TicketStatus
	err := row.Scan(&status)
	return status, err
}

const markExpiredReservations = `-- name: MarkExpiredReservations :exec
UPDATE reservation
SET status = 'CANCELED-BY-TIME'
WHERE (created_at + duration_time) < now()
  AND status != 'RESERVED'
`

func (q *Queries) MarkExpiredReservations(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, markExpiredReservations)
	return err
}

const updateReservation = `-- name: UpdateReservation :one
UPDATE "reservation"
SET 
    "status" = $1
WHERE id = $2
RETURNING 
    id, 
    user_id, 
    ticket_id, 
    payment_id, 
    status, 
    EXTRACT(EPOCH FROM duration_time)::bigint as duration_time_seconds,
    created_at
`

type UpdateReservationParams struct {
	Status TicketStatus `json:"status"`
	ID     int64        `json:"id"`
}

type UpdateReservationRow struct {
	ID                  int64        `json:"id"`
	UserID              int64        `json:"user_id"`
	TicketID            int64        `json:"ticket_id"`
	PaymentID           int64        `json:"payment_id"`
	Status              TicketStatus `json:"status"`
	DurationTimeSeconds int64        `json:"duration_time_seconds"`
	CreatedAt           time.Time    `json:"created_at"`
}

func (q *Queries) UpdateReservation(ctx context.Context, arg UpdateReservationParams) (UpdateReservationRow, error) {
	row := q.db.QueryRowContext(ctx, updateReservation, arg.Status, arg.ID)
	var i UpdateReservationRow
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.TicketID,
		&i.PaymentID,
		&i.Status,
		&i.DurationTimeSeconds,
		&i.CreatedAt,
	)
	return i, err
}
