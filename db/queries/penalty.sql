-- name: GetTicketPenalties :many
SELECT p.* 
FROM penalty p
JOIN ticket t ON p.vehicle_id = t.vehicle_id
WHERE t.id = $1;