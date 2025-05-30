-- db/migrate/000005_default_penalty.up.sql

-- Insert default penalty rules for bus vehicles
INSERT INTO "penalty" ("vehicle_id", "penalty_text", "befor_day", "after_day") VALUES
-- Iran Peyma buses (vehicles 1-3)
(1, 'Bus cancellation penalty: 10% before departure, 25% after departure', 10, 25),
(2, 'Bus cancellation penalty: 10% before departure, 25% after departure', 10, 25),
(3, 'VIP Bus cancellation penalty: 15% before departure, 35% after departure', 15, 35),

-- Hamrah Safar buses (vehicles 4-6)
(4, 'Bus cancellation penalty: 12% before departure, 30% after departure', 12, 30),
(5, 'Bus cancellation penalty: 12% before departure, 30% after departure', 12, 30),
(6, 'VIP Bus cancellation penalty: 18% before departure, 40% after departure', 18, 40),

-- Royal Safar buses (vehicles 7-8)
(7, 'Premium bus cancellation penalty: 15% before departure, 35% after departure', 15, 35),
(8, 'VIP Bus cancellation penalty: 20% before departure, 45% after departure', 20, 45),

-- Hamsafar buses (vehicles 9-10)
(9, 'Bus cancellation penalty: 10% before departure, 25% after departure', 10, 25),
(10, 'VIP Bus cancellation penalty: 18% before departure, 40% after departure', 18, 40);

-- Insert default penalty rules for train vehicles
INSERT INTO "penalty" ("vehicle_id", "penalty_text", "befor_day", "after_day") VALUES
-- Raja Rail trains (vehicles 11-13)
(11, 'Train cancellation penalty: 20% before departure, 50% after departure', 20, 50),
(12, 'Train cancellation penalty: 20% before departure, 50% after departure', 20, 50),
(13, 'Standard train cancellation penalty: 15% before departure, 40% after departure', 15, 40),

-- Fadak Rail trains (vehicles 14-15)
(14, 'Train cancellation penalty: 18% before departure, 45% after departure', 18, 45),
(15, 'Premium train cancellation penalty: 25% before departure, 60% after departure', 25, 60),

-- Safir Rail trains (vehicle 16)
(16, 'Train cancellation penalty: 15% before departure, 40% after departure', 15, 40),

-- Persian Rail trains (vehicle 17)
(17, 'Express train cancellation penalty: 22% before departure, 55% after departure', 22, 55);

-- Insert default penalty rules for airplane vehicles
INSERT INTO "penalty" ("vehicle_id", "penalty_text", "befor_day", "after_day") VALUES
-- Iran Air planes (vehicles 18-20)
(18, 'Economy flight cancellation penalty: 30% before departure, 75% after departure', 30, 75),
(19, 'Economy flight cancellation penalty: 30% before departure, 75% after departure', 30, 75),
(20, 'Business class cancellation penalty: 40% before departure, 85% after departure', 40, 85),

-- Mahan Air planes (vehicles 21-22)
(21, 'Economy flight cancellation penalty: 35% before departure, 80% after departure', 35, 80),
(22, 'Business class cancellation penalty: 45% before departure, 90% after departure', 45, 90),

-- Aseman Airlines planes (vehicles 23-24)
(23, 'Economy flight cancellation penalty: 25% before departure, 70% after departure', 25, 70),
(24, 'Premium economy cancellation penalty: 35% before departure, 80% after departure', 35, 80),

-- Iran Air Tours planes (vehicle 25)
(25, 'Economy flight cancellation penalty: 28% before departure, 72% after departure', 28, 72),

-- Kish Air planes (vehicle 26)
(26, 'Premium economy cancellation penalty: 38% before departure, 82% after departure', 38, 82);

-- Insert additional penalty rules for different time periods
-- Early cancellation penalties (more than 7 days before departure)
INSERT INTO "penalty" ("vehicle_id", "penalty_text", "befor_day", "after_day") VALUES
-- Sample early cancellation penalties for some vehicles
(1, 'Early bus cancellation (7+ days): 5% penalty, standard rates apply after', 5, 25),
(11, 'Early train cancellation (7+ days): 10% penalty, standard rates apply after', 10, 50),
(18, 'Early flight cancellation (7+ days): 15% penalty, standard rates apply after', 15, 75),

-- Late booking penalties (less than 24 hours before departure)
(3, 'Last-minute VIP bus booking: 20% surcharge before, 50% penalty after', 20, 50),
(15, 'Last-minute premium train booking: 30% surcharge before, 70% penalty after', 30, 70),
(22, 'Last-minute business flight booking: 50% surcharge before, 100% penalty after', 50, 100);

-- Insert seasonal penalty adjustments
INSERT INTO "penalty" ("vehicle_id", "penalty_text", "befor_day", "after_day") VALUES
-- Holiday season penalties
(2, 'Holiday season bus penalty: 15% before departure, 35% after departure', 15, 35),
(12, 'Holiday season train penalty: 25% before departure, 60% after departure', 25, 60),
(19, 'Holiday season flight penalty: 40% before departure, 85% after departure', 40, 85);

-- Create index for better performance on penalty queries
CREATE INDEX ON "penalty" ("vehicle_id");
CREATE INDEX ON "penalty" ("vehicle_id", "befor_day", "after_day");