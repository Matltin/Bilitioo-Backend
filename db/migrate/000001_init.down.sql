-- First drop notification and verification related tables (most dependent)
DROP TABLE IF EXISTS "notification_log";
DROP TABLE IF EXISTS "send_verification_code";
DROP TABLE IF EXISTS "send_email_sms";

-- Drop user activity tables
ALTER TABLE "user_activity" DROP CONSTRAINT IF EXISTS "user_activity_user_id_fkey";
ALTER TABLE "user_activity" DROP CONSTRAINT IF EXISTS "user_activity_route_id_fkey";
DROP TABLE IF EXISTS "user_activity";

-- Drop report tables
ALTER TABLE "report" DROP CONSTRAINT IF EXISTS "report_reservation_id_fkey";
ALTER TABLE "report" DROP CONSTRAINT IF EXISTS "report_admin_id_fkey";
ALTER TABLE "report" DROP CONSTRAINT IF EXISTS "report_user_id_fkey";
DROP TABLE IF EXISTS "report";

-- Drop change reservation tables
ALTER TABLE "change_reservation" DROP CONSTRAINT IF EXISTS "change_reservation_reservation_id_fkey";
ALTER TABLE "change_reservation" DROP CONSTRAINT IF EXISTS "change_reservation_admin_id_fkey";
ALTER TABLE "change_reservation" DROP CONSTRAINT IF EXISTS "change_reservation_user_id_fkey";
DROP TABLE IF EXISTS "change_reservation";

-- Drop reservation tables
ALTER TABLE "reservation" DROP CONSTRAINT IF EXISTS "reservation_user_id_fkey";
ALTER TABLE "reservation" DROP CONSTRAINT IF EXISTS "reservation_ticket_id_fkey";
ALTER TABLE "reservation" DROP CONSTRAINT IF EXISTS "reservation_payment_id_fkey";
DROP TABLE IF EXISTS "reservation";

-- Drop payment tables
DROP TABLE IF EXISTS "payment";

-- Drop penalty tables
ALTER TABLE "penalty" DROP CONSTRAINT IF EXISTS "penalty_vehicle_id_fkey";
DROP TABLE IF EXISTS "penalty";

-- Drop seat-related tables (depends on vehicle)
ALTER TABLE "bus_seat" DROP CONSTRAINT IF EXISTS "bus_seat_seat_id_fkey";
ALTER TABLE "train_seat" DROP CONSTRAINT IF EXISTS "train_seat_seat_id_fkey";
ALTER TABLE "airplane_seat" DROP CONSTRAINT IF EXISTS "airplane_seat_seat_id_fkey";
DROP TABLE IF EXISTS "bus_seat";
DROP TABLE IF EXISTS "train_seat";
DROP TABLE IF EXISTS "airplane_seat";

ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS "ticket_seat_id_fkey";
DROP TABLE IF EXISTS "seat";

-- Drop ticket tables (depends on route and vehicle)
ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS "ticket_vehicle_id_fkey";
ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS "ticket_route_id_fkey";
DROP TABLE IF EXISTS "ticket";

-- Drop vehicle-specific tables
ALTER TABLE "bus" DROP CONSTRAINT IF EXISTS "bus_vehicle_id_fkey";
ALTER TABLE "train" DROP CONSTRAINT IF EXISTS "train_vehicle_id_fkey";
ALTER TABLE "airplane" DROP CONSTRAINT IF EXISTS "airplane_vehicle_id_fkey";
DROP TABLE IF EXISTS "bus";
DROP TABLE IF EXISTS "train";
DROP TABLE IF EXISTS "airplane";

-- Drop vehicle tables (depends on company)
ALTER TABLE "vehicle" DROP CONSTRAINT IF EXISTS "vehicle_company_id_fkey";
DROP TABLE IF EXISTS "vehicle";

-- Drop company tables
DROP TABLE IF EXISTS "company";

-- Drop route tables (depends on city and terminal)
ALTER TABLE "route" DROP CONSTRAINT IF EXISTS "route_origin_city_id_fkey";
ALTER TABLE "route" DROP CONSTRAINT IF EXISTS "route_destination_city_id_fkey";
ALTER TABLE "route" DROP CONSTRAINT IF EXISTS "route_origin_terminal_id_fkey";
ALTER TABLE "route" DROP CONSTRAINT IF EXISTS "route_destination_terminal_id_fkey";
DROP TABLE IF EXISTS "route";
DROP TABLE IF EXISTS "terminal";

-- Drop profile tables (depends on user and city)
ALTER TABLE "profile" DROP CONSTRAINT IF EXISTS "profile_user_id_fkey";
ALTER TABLE "profile" DROP CONSTRAINT IF EXISTS "profile_city_id_fkey";
DROP TABLE IF EXISTS "profile";

-- Drop city tables
DROP TABLE IF EXISTS "city";

-- Drop user tables
ALTER TABLE "user" DROP CONSTRAINT IF EXISTS "email_or_phone_required";
DROP TABLE IF EXISTS "user";



-- Drop all indexes
DROP INDEX IF EXISTS "user_email_idx";
DROP INDEX IF EXISTS "user_phone_number_idx";
DROP INDEX IF EXISTS "profile_user_id_idx";
DROP INDEX IF EXISTS "city_province_idx";
DROP INDEX IF EXISTS "city_county_idx";
DROP INDEX IF EXISTS "city_province_county_idx";
DROP INDEX IF EXISTS "ticket_route_id_idx";
DROP INDEX IF EXISTS "ticket_departure_time_idx";
DROP INDEX IF EXISTS "ticket_route_id_departure_time_vehicle_id_idx";
DROP INDEX IF EXISTS "vehicle_company_id_idx";
DROP INDEX IF EXISTS "route_origin_city_id_idx";
DROP INDEX IF EXISTS "route_destination_city_id_idx";
DROP INDEX IF EXISTS "route_origin_terminal_id_idx";
DROP INDEX IF EXISTS "route_origin_city_id_destination_city_id_idx";
DROP INDEX IF EXISTS "route_destination_terminal_id_idx";
DROP INDEX IF EXISTS "route_origin_terminal_id_destination_terminal_id_idx";
DROP INDEX IF EXISTS "seat_vehicle_id_idx";
DROP INDEX IF EXISTS "reservation_user_id_idx";
DROP INDEX IF EXISTS "reservation_ticket_id_idx";
DROP INDEX IF EXISTS "reservation_user_id_ticket_id_idx";
DROP INDEX IF EXISTS "report_user_id_idx";
DROP INDEX IF EXISTS "report_admin_id_idx";
DROP INDEX IF EXISTS "report_user_id_admin_id_idx";

-- Drop all types
DROP TYPE IF EXISTS "activity_status";
DROP TYPE IF EXISTS "notification_log_status";
DROP TYPE IF EXISTS "request_type";
DROP TYPE IF EXISTS "ticket_status";
DROP TYPE IF EXISTS "payment_type";
DROP TYPE IF EXISTS "payment_status";
DROP TYPE IF EXISTS "check_reservation_ticket_status";
DROP TYPE IF EXISTS "vehicle_type";
DROP TYPE IF EXISTS "flight_class";
DROP TYPE IF EXISTS "role";
DROP TYPE IF EXISTS "user_status";
