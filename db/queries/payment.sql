-- name: CreatePayment :one
INSERT INTO "payment" (
    "from_account",
    "to_account",
    "amount"
) VALUES (
    $1, $2, $3
) RETURNING *;