-- name: GetTicketDetails :one
SELECT 
    oc.province AS "origin",
    dc.province AS "destination",
    t.departure_time,
    t.arrival_time,
    t.amount,
    v.capacity,
    v.vehicle_type,
    v.feature,
    b."VIP",
    b.bed_chair,
    tr.rank,
    tr.have_compartment,
    a.flight_class,
    a.name AS airplane_name

FROM "ticket" t
INNER JOIN "vehicle" v ON t.vehicle_id = v.id
LEFT JOIN "bus" b ON b.vehicle_id = v.id 
LEFT JOIN "train" tr ON tr.vehicle_id = v.id 
LEFT JOIN "airplane" a ON a.vehicle_id = v.id 
INNER JOIN "route" ro ON t.route_id = ro.id
INNER JOIN "city" oc ON oc.id = ro.origin_city_id
INNER JOIN "city" dc ON dc.id = ro.destination_city_id
WHERE t.id = $1;

-- name: SearchTickets :many
SELECT 
    t.id,
    t.vehicle_id,
    t.seat_id,
    t.vehicle_type,
    t.route_id,
    t.amount,
    to_char(t.departure_time, 'YYYY-MM-DD HH24:MI') as departure_time,
    to_char(t.arrival_time, 'YYYY-MM-DD HH24:MI') as arrival_time,
    t.count_stand,
    t.status
FROM ticket t
JOIN route r ON t.route_id = r.id
WHERE
    ($1::bigint IS NULL OR r.origin_city_id = $1)
    AND ($2::bigint IS NULL OR r.destination_city_id = $2)
    AND ($3::date IS NULL OR t.departure_time::date = $3::date)
    AND ($4::vehicle_type IS NULL OR t.vehicle_type = $4)
    AND t.status = 'NOT_RESERVED'
ORDER BY t.departure_time ASC
LIMIT 50;
