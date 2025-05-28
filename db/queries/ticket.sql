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
