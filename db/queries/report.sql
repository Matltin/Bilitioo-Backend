-- name: GetReports :many
SELECT * FROM "report";

-- name: AnswerReport :one
UPDATE "report"
SET response_text = $1
WHERE id = $2
RETURNING *;