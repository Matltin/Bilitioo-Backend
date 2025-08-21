INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  22, -- vehicle_id
  s.id, -- seat_id
  'AIRPLANE', -- vehicle_type
  30, -- route_id
  3000000, -- amount (مثلاً 3,000,000 Rials)
  NOW() + INTERVAL '1 day' + INTERVAL '15 hour', -- departure_time (فردا ساعت 15:00)
  NOW() + INTERVAL '1 day' + INTERVAL '16 hour' + INTERVAL '40 minute', -- arrival_time (فردا 16:40)
  0, -- count_stand
  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 22
LIMIT 1; -- فقط یک بلیط ساخته بشه