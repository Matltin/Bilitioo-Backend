-- name: GetCities :many
SELECT 
    "province", 
    "county"
FROM "city";

-- name: SearchTicketsByCities :many
SELECT 
    t.id AS ticket_id,
    t.departure_time,
    t.arrival_time,
    t.amount,
    t.status,
    r.id AS route_id,
    r.origin_city_id,
    r.destination_city_id,
    c1.province AS origin_province,
    c1.county AS origin_county,
    c2.province AS destination_province,
    c2.county AS destination_county,
    v.vehicle_type,
    v.capacity,
    comp.name AS company_name
FROM ticket t
JOIN route r ON t.route_id = r.id
JOIN city c1 ON r.origin_city_id = c1.id
JOIN city c2 ON r.destination_city_id = c2.id
JOIN vehicle v ON t.vehicle_id = v.id
JOIN company comp ON v.company_id = comp.id
WHERE r.origin_city_id = $1 AND r.destination_city_id = $2 AND v.vehicle_type = $3
ORDER BY t.departure_time ASC;