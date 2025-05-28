-- name: CreatePayment :one
INSERT INTO "payment" (
    "from_account",
    "to_account",
    "amount"
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdatePayment :one
UPDATE "payment"
SET 
    "type" = $1,
    "status" = $2
WHERE id = $3
RETURNING *;