-- name: CreateChangeReservation :one
INSERT INTO "change_reservation" (
    "reservation_id",
    "admin_id",
    "user_id",
    "from_status",
    "to_status"
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;