-- Tickets for buses (only adding tickets for vehicles 1, 3, 5, and 8)
-- Iran Peyma regular bus (vehicle 1): Tehran to Isfahan
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  1, -- vehicle_id
  s.id, -- seat_id
  'BUS', -- vehicle_type
  4, -- route_id (Tehran to Isfahan)
  650000, -- amount (650,000 Rials = ~$15)
  NOW() + INTERVAL '1 day' + INTERVAL '8 hour', -- departure_time (tomorrow at 8:00 AM)
  NOW() + INTERVAL '1 day' + INTERVAL '14 hour', -- arrival_time (tomorrow at 2:00 PM)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 1 AND s.seat_number <= 35; -- Only first 40 seats available

-- Iran Peyma VIP bus (vehicle 3): Tehran to Shiraz
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  3, -- vehicle_id
  s.id, -- seat_id
  'BUS', -- vehicle_type
  7, -- route_id (Tehran to Shiraz)
  1200000, -- amount (1,200,000 Rials = ~$28)
  NOW() + INTERVAL '2 day' + INTERVAL '20 hour', -- departure_time (day after tomorrow at 8:00 PM)
  NOW() + INTERVAL '3 day' + INTERVAL '6 hour', -- arrival_time (3 days from now at 6:00 AM)
  0, -- count_stand
  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 3;

-- Hamrah Safar regular bus (vehicle 5): Mashhad to Tehran
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  5, -- vehicle_id
  s.id, -- seat_id
  'BUS', -- vehicle_type
  31, -- route_id (Mashhad to Tehran)
  750000, -- amount (750,000 Rials = ~$18)
  NOW() + INTERVAL '1 day' + INTERVAL '22 hour', -- departure_time (tomorrow at 10:00 PM)
  NOW() + INTERVAL '2 day' + INTERVAL '10 hour', -- arrival_time (day after tomorrow at 10:00 AM)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 5;

-- Royal Safar VIP bus (vehicle 8): Isfahan to Shiraz
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  8, -- vehicle_id
  s.id, -- seat_id
  'BUS', -- vehicle_type
  26, -- route_id (Isfahan to Shiraz)
  950000, -- amount (950,000 Rials = ~$22)
  NOW() + INTERVAL '3 day' + INTERVAL '14 hour', -- departure_time (3 days from now at 2:00 PM)
  NOW() + INTERVAL '3 day' + INTERVAL '20 hour', -- arrival_time (3 days from now at 8:00 PM)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 8;

-- Tickets for trains (only adding tickets for vehicles 11, 15, and 17)
-- Raja Rail train (vehicle 11): Tehran to Mashhad
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  11, -- vehicle_id
  s.id, -- seat_id
  'TRAIN', -- vehicle_type
  3, -- route_id (Tehran to Mashhad by train)
  1500000, -- amount (1,500,000 Rials = ~$35)
  NOW() + INTERVAL '2 day' + INTERVAL '16 hour', -- departure_time (day after tomorrow at 4:00 PM)
  NOW() + INTERVAL '3 day' + INTERVAL '8 hour', -- arrival_time (3 days from now at 8:00 AM)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 11 AND s.seat_number <= 350; -- Only 350 seats available out of 400

-- Fadak Rail premium train (vehicle 15): Tehran to Isfahan
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  15, -- vehicle_id
  s.id, -- seat_id
  'TRAIN', -- vehicle_type
  6, -- route_id (Tehran to Isfahan by train)
  1800000, -- amount (1,800,000 Rials = ~$42)
  NOW() + INTERVAL '4 day' + INTERVAL '9 hour', -- departure_time (4 days from now at 9:00 AM)
  NOW() + INTERVAL '4 day' + INTERVAL '15 hour', -- arrival_time (4 days from now at 3:00 PM)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 15;

-- Persian Rail train (vehicle 17): Isfahan to Tehran
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  17, -- vehicle_id
  s.id, -- seat_id
  'TRAIN', -- vehicle_type
  32, -- route_id (Isfahan to Tehran)
  1650000, -- amount (1,650,000 Rials = ~$38)
  NOW() + INTERVAL '3 day' + INTERVAL '18 hour', -- departure_time (3 days from now at 6:00 PM)
  NOW() + INTERVAL '4 day' + INTERVAL '00 hour', -- arrival_time (4 days from now at midnight)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 17;

-- Tickets for airplanes (only adding tickets for vehicles 18, 22, and 26)
-- Iran Air economy plane (vehicle 18): Tehran to Mashhad
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  18, -- vehicle_id
  s.id, -- seat_id
  'AIRPLANE', -- vehicle_type
  2, -- route_id (Tehran to Mashhad by air)
  2500000, -- amount (2,500,000 Rials = ~$58)
  NOW() + INTERVAL '1 day' + INTERVAL '10 hour', -- departure_time (tomorrow at 10:00 AM)
  NOW() + INTERVAL '1 day' + INTERVAL '11 hour' + INTERVAL '30 minute', -- arrival_time (tomorrow at 11:30 AM)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 18;

-- Mahan Air business plane (vehicle 22): Tehran to Shiraz
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  22, -- vehicle_id
  s.id, -- seat_id
  'AIRPLANE', -- vehicle_type
  8, -- route_id (Tehran to Shiraz by air)
  3800000, -- amount (3,800,000 Rials = ~$88)
  NOW() + INTERVAL '2 day' + INTERVAL '15 hour', -- departure_time (day after tomorrow at 3:00 PM)
  NOW() + INTERVAL '2 day' + INTERVAL '16 hour' + INTERVAL '20 minute', -- arrival_time (day after tomorrow at 4:20 PM)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 22;

-- Kish Air premium economy plane (vehicle 26): Isfahan to Tehran
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
SELECT 
  26, -- vehicle_id
  s.id, -- seat_id
  'AIRPLANE', -- vehicle_type
  32, -- route_id (Isfahan to Tehran)
  3200000, -- amount (3,200,000 Rials = ~$74)
  NOW() + INTERVAL '5 day' + INTERVAL '7 hour', -- departure_time (5 days from now at 7:00 AM)
  NOW() + INTERVAL '5 day' + INTERVAL '8 hour' + INTERVAL '10 minute', -- arrival_time (5 days from now at 8:10 AM)
  0, -- count_stand

  'NOT_RESERVED' -- status
FROM "seat" s
WHERE s.vehicle_id = 26;

-- Create some "standing" tickets for one bus route (Tehran to Isfahan)
INSERT INTO "ticket" (
  "vehicle_id", "seat_id", "vehicle_type", "route_id", "amount", 
  "departure_time", "arrival_time", "count_stand", "status"
)
VALUES (
  1, -- vehicle_id (Iran Peyma regular bus)
  (SELECT MIN(id) FROM seat WHERE vehicle_id = 1), -- Using first seat as reference
  'BUS', -- vehicle_type
  4, -- route_id (Tehran to Isfahan)
  450000, -- amount (450,000 Rials = ~$10)
  NOW() + INTERVAL '1 day' + INTERVAL '8 hour', -- departure_time (same as the seated bus)
  NOW() + INTERVAL '1 day' + INTERVAL '14 hour', -- arrival_time
  5, -- count_stand (5 standing tickets available)
  'NOT_RESERVED' -- status
);

