  -- First, let's insert transportation companies for each vehicle type
INSERT INTO "company" ("name", "address") VALUES
-- Bus companies
('Iran Peyma', 'Tehran, Azadi Ave'),
('Hamrah Safar', 'Tehran, Enqelab St'),
('Royal Safar', 'Isfahan, Chaharbagh Ave'),
('Hamsafar', 'Mashhad, Imam Reza Blvd'),
-- Train companies
('Raja Rail Transportation', 'Tehran, Railway Sq'),
('Fadak Rail', 'Tehran, Jomhouri St'),
('Safir Rail', 'Mashhad, Railway Sq'),
('Persian Rail', 'Isfahan, Railway Sq'),
-- Airline companies
('Iran Air', 'Tehran, Mehrabad Airport'),
('Mahan Air', 'Tehran, Azadi St'),
('Aseman Airlines', 'Tehran, Vali Asr St'),
('Iran Air Tours', 'Mashhad, Airport Blvd'),
('Kish Air', 'Tehran, Vanak Sq');

-- Now, let's create vehicles for each company
-- First for bus companies
INSERT INTO "vehicle" ("company_id", "capacity", "vehicle_type","feature")
VALUES
-- Iran Peyma buses
(1, 44, 'BUS','{"air_conditioning": true, "wifi": true, "restroom": true}'),
(1, 44, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true}'),
(1, 25, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true, "premium_seats": true}'),
-- Hamrah Safar buses
(2, 44, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true}'),
(2, 44, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true}'),
(2, 25, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true, "premium_seats": true}'),
-- Royal Safar buses
(3, 44, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true}'),
(3, 25, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true, "premium_seats": true}'),
-- Hamsafar buses
(4, 44, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true}'),
(4, 25, 'BUS', '{"air_conditioning": true, "wifi": true, "restroom": true, "premium_seats": true}');

-- Now for train companies
INSERT INTO "vehicle" ("company_id", "capacity", "vehicle_type", "feature")
VALUES
-- Raja Rail trains
(5, 400, 'TRAIN', '{"dining_car": true, "sleeper_cabins": true, "wifi": true}'),
(5, 400, 'TRAIN', '{"dining_car": true, "sleeper_cabins": true, "wifi": true}'),
(5, 350, 'TRAIN', '{"dining_car": true, "sleeper_cabins": false, "wifi": true}'),
-- Fadak Rail trains
(6, 350, 'TRAIN', '{"dining_car": true, "sleeper_cabins": false, "wifi": true}'),
(6, 450, 'TRAIN', '{"dining_car": true, "sleeper_cabins": true, "wifi": true, "premium_service": true}'),
-- Safir Rail trains
(7, 300, 'TRAIN', '{"dining_car": true, "sleeper_cabins": false, "wifi": true}'),
-- Persian Rail trains
(8, 400, 'TRAIN', '{"dining_car": true, "sleeper_cabins": true, "wifi": true}');

-- Airplanes
INSERT INTO "vehicle" ("company_id", "capacity", "vehicle_type", "feature")
VALUES
-- Iran Air planes
(9, 150, 'AIRPLANE', '{"meal_service": true, "entertainment": true}'),
(9, 180, 'AIRPLANE', '{"meal_service": true, "entertainment": true}'),
(9, 220, 'AIRPLANE', '{"meal_service": true, "entertainment": true, "premium_cabin": true}'),
-- Mahan Air planes
(10, 160, 'AIRPLANE', '{"meal_service": true, "entertainment": true}'),
(10, 200, 'AIRPLANE', '{"meal_service": true, "entertainment": true, "premium_cabin": true}'),
-- Aseman Airlines planes
(11, 120, 'AIRPLANE', '{"meal_service": true, "entertainment": false}'),
(11, 140, 'AIRPLANE', '{"meal_service": true, "entertainment": true}'),
-- Iran Air Tours planes
(12, 130, 'AIRPLANE', '{"meal_service": true, "entertainment": false}'),
-- Kish Air planes
(13, 150, 'AIRPLANE', '{"meal_service": true, "entertainment": true}');

-- Now, let's add specific details for each vehicle type
-- Buses
INSERT INTO "bus" ("vehicle_id", "VIP", "bed_chair")
VALUES
(1, false, false),
(2, false, false),
(3, true, true),
(4, false, false),
(5, false, false),
(6, true, true),
(7, false, false),
(8, true, true),
(9, false, false),
(10, true, true);

-- Trains
INSERT INTO "train" ("vehicle_id", "rank", "have_compartment")
VALUES
(11, 4, true),
(12, 4, true),
(13, 3, false),
(14, 3, false),
(15, 5, true),
(16, 3, false),
(17, 4, true);

-- Airplanes
INSERT INTO "airplane" ("vehicle_id", "flight_class", "name")
VALUES
(18, 'ECONOMY', 'Airbus A320'),
(19, 'ECONOMY', 'Boeing 737'),
(20, 'BUSINESS', 'Airbus A330'),
(21, 'ECONOMY', 'Boeing 737'),
(22, 'BUSINESS', 'Airbus A320'),
(23, 'ECONOMY', 'Fokker 100'),
(24, 'PREMIUM-ECONOMY', 'Boeing 737'),
(25, 'ECONOMY', 'MD-80'),
(26, 'PREMIUM-ECONOMY', 'Boeing 737');

-- Insert seat 
DO $$
DECLARE
  rec RECORD;
  i INT;
BEGIN
  FOR rec IN (
    SELECT v.id, v.capacity, v.vehicle_type AS vehicle_type
    FROM vehicle v
  ) LOOP
    FOR i IN 1..rec.capacity LOOP
      INSERT INTO seat (vehicle_id, vehicle_type, seat_number)
      VALUES (rec.id, rec.vehicle_type, i);
    END LOOP;
  END LOOP;
END;
$$;

INSERT INTO bus_seat (seat_id)
SELECT s.id
FROM seat s
JOIN vehicle v ON s.vehicle_id = v.id
WHERE s.vehicle_type = 'BUS';

INSERT INTO bus_seat (seat_id)
SELECT s.id
FROM seat s
JOIN vehicle v ON s.vehicle_id = v.id
WHERE s.vehicle_type = 'AIRPLANE';

INSERT INTO train_seat (seat_id, salon, have_compartment, cuope_number)
SELECT 
	s.id,
	(s.seat_number - 1) / 50 + 1 AS salon,
 	t.have_compartment,
  	((s.seat_number - 1) % 50) / 6 + 1 AS cuope_number
FROM seat s
JOIN vehicle v ON s.vehicle_id = v.id
JOIN train t ON t.vehicle_id = v.id;
