-- name: GetReports :many
SELECT * FROM "report";

-- name: AnswerReport :one
UPDATE "report"
SET 
    response_text = $1,
    admin_id = $2
WHERE id = $3
RETURNING *;

-- name: CreateReport :one
INSERT INTO "report" (
    "reservation_id",
    "user_id",
    "request_type",
    "request_text",
    "response_text"
)
VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;