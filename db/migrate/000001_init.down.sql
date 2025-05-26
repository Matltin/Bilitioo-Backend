ALTER TABLE "profile" DROP CONSTRAINT "profile_user_id_fkey";
ALTER TABLE "profile" DROP CONSTRAINT "profile_city_id_fkey";

DROP INDEX IF EXISTS "user_email_idx";
DROP INDEX IF EXISTS "user_phone_number_idx";
DROP INDEX IF EXISTS "profile_user_id_idx";
DROP INDEX IF EXISTS "city_province_idx";
DROP INDEX IF EXISTS "city_county_idx";
DROP INDEX IF EXISTS "city_province_county_idx";


ALTER TABLE "user" DROP CONSTRAINT IF EXISTS phone_number_check;
ALTER TABLE "user" DROP CONSTRAINT IF EXISTS email_check;
ALTER TABLE "user" DROP CONSTRAINT IF EXISTS email_or_phone_required;

DROP TABLE IF EXISTS "profile";
DROP TABLE IF EXISTS "city";
DROP TABLE IF EXISTS "user";

DROP TYPE IF EXISTS "role";
DROP TYPE IF EXISTS "user_status";

ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS "ticket_vehicle_id_fkey";
ALTER TABLE "vehicle" DROP CONSTRAINT IF EXISTS "vehicle_company_id_fkey";
ALTER TABLE "bus" DROP CONSTRAINT IF EXISTS "bus_vehicle_id_fkey";
ALTER TABLE "train" DROP CONSTRAINT IF EXISTS "train_vehicle_id_fkey";
ALTER TABLE "airplane" DROP CONSTRAINT IF EXISTS "airplane_vehicle_id_fkey";

DROP INDEX IF EXISTS "ticket_route_id_idx";
DROP INDEX IF EXISTS "ticket_departure_time_idx";
DROP INDEX IF EXISTS "ticket_route_id_departure_time_vehicle_id_idx";
DROP INDEX IF EXISTS "vehicle_company_id_idx";

ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS amount_validation;
ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS time_validation;
ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS count_stand_validation;

ALTER TABLE "vehicle" DROP CONSTRAINT IF EXISTS capacity_validation;

ALTER TABLE "train" DROP CONSTRAINT IF EXISTS rank_validation;

DROP TABLE IF EXISTS "ticket";
DROP TABLE IF EXISTS "vehicle";
DROP TABLE IF EXISTS "company";
DROP TABLE IF EXISTS "bus";
DROP TABLE IF EXISTS "train";
DROP TABLE IF EXISTS "airplane";

DROP TYPE IF EXISTS "check_reservation_ticket_status";
DROP TYPE IF EXISTS "vehicle_type";
DROP TYPE IF EXISTS "flight_class";

ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS "ticket_route_id_fkey";
ALTER TABLE "route" DROP CONSTRAINT IF EXISTS "route_origin_city_id_fkey";
ALTER TABLE "route" DROP CONSTRAINT IF EXISTS "route_destination_city_id_fkey";
ALTER TABLE "route" DROP CONSTRAINT IF EXISTS "route_origin_terminal_id_fkey";
ALTER TABLE "route" DROP CONSTRAINT IF EXISTS "route_destination_terminal_id_fkey";

DROP INDEX IF EXISTS "route_origin_city_id_idx";
DROP INDEX IF EXISTS "route_destination_city_id_idx";
DROP INDEX IF EXISTS "route_origin_terminal_id_idx";
DROP INDEX IF EXISTS "route_origin_city_id_destination_city_id_idx";
DROP INDEX IF EXISTS "route_destination_terminal_id_idx";
DROP INDEX IF EXISTS "route_origin_terminal_id_destination_terminal_id_idx";

ALTER TABLE "route" DROP CONSTRAINT IF EXISTS distance_validation;

DROP TABLE IF EXISTS "route";
DROP TABLE IF EXISTS "terminal";

ALTER TABLE "ticket" DROP CONSTRAINT IF EXISTS "ticket_seat_id_fkey";
ALTER TABLE "seat" DROP CONSTRAINT IF EXISTS "seat_vehicle_id_fkey";
ALTER TABLE "bus_seat" DROP CONSTRAINT IF EXISTS "bus_seat_seat_id_fkey";
ALTER TABLE "train_seat" DROP CONSTRAINT IF EXISTS "train_seat_seat_id_fkey";
ALTER TABLE "airplane_seat" DROP CONSTRAINT IF EXISTS "airplane_seat_seat_id_fkey";

DROP INDEX IF EXISTS "seat_vehicle_id_idx";

ALTER TABLE "seat" DROP CONSTRAINT IF EXISTS seat_number_validation;

ALTER TABLE "train_seat" DROP CONSTRAINT IF EXISTS salon_validation;
ALTER TABLE "train_seat" DROP CONSTRAINT IF EXISTS coupe_number_validation;

DROP TABLE IF EXISTS "bus_seat";
DROP TABLE IF EXISTS "train_seat";
DROP TABLE IF EXISTS "airplane_seat";
DROP TABLE IF EXISTS "seat";


ALTER TABLE "reservation" DROP CONSTRAINT IF EXISTS "reservation_user_id_fkey";
ALTER TABLE "reservation" DROP CONSTRAINT IF EXISTS "reservation_ticket_id_fkey";
ALTER TABLE "reservation" DROP CONSTRAINT IF EXISTS "reservation_payment_id_fkey";
ALTER TABLE "change_reservation" DROP CONSTRAINT IF EXISTS "change_reservation_reservation_id_fkey";
ALTER TABLE "change_reservation" DROP CONSTRAINT IF EXISTS "change_reservation_admin_id_fkey";
ALTER TABLE "change_reservation" DROP CONSTRAINT IF EXISTS "change_reservation_user_id_fkey";
ALTER TABLE "report" DROP CONSTRAINT "report_admin_id_fkey";
ALTER TABLE "report" DROP CONSTRAINT "report_user_id_fkey";

DROP INDEX IF EXISTS "reservation_user_id_idx";
DROP INDEX IF EXISTS "reservation_ticket_id_idx";
DROP INDEX IF EXISTS "reservation_user_id_ticket_id_idx";
DROP INDEX IF EXISTS "report_user_id_idx";
DROP INDEX IF EXISTS "report_admin_id_idx";
DROP INDEX IF EXISTS "report_user_id_admin_id_idx";

ALTER TABLE "payment" DROP CONSTRAINT IF EXISTS amount_payment_validation;

ALTER TABLE "reservation" DROP CONSTRAINT IF EXISTS duration_time_reservation_validation;

DROP TABLE IF EXISTS "report";
DROP TABLE IF EXISTS "change_reservation";
DROP TABLE IF EXISTS "reservation";
DROP TABLE IF EXISTS "payment";


DROP TYPE IF EXISTS "request_type";
DROP TYPE IF EXISTS "ticket_status";
DROP TYPE IF EXISTS "payment_type";
DROP TYPE IF EXISTS "payment_status";

ALTER TABLE "penalty" DROP CONSTRAINT IF EXISTS "penalty_vehicle_id_fkey";

DROP TABLE IF EXISTS "penalty";

ALTER TABLE "user_activity" DROP CONSTRAINT IF EXISTS "user_activity_user_id_fkey";
ALTER TABLE "user_activity" DROP CONSTRAINT IF EXISTS "user_activity_route_id_fkey";
ALTER TABLE "notification_log" DROP CONSTRAINT IF EXISTS "notification_log_user_id_fkey";
ALTER TABLE "notification_log" DROP CONSTRAINT IF EXISTS "notification_log_send_email_sms_id_fkey";
ALTER TABLE "notification_log" DROP CONSTRAINT IF EXISTS "notification_log_user_activity_id_fkey";
ALTER TABLE "send_verification_code" DROP CONSTRAINT IF EXISTS "send_verification_code_user_id_fkey";
ALTER TABLE "send_verification_code" DROP CONSTRAINT IF EXISTS "send_verification_code_send_email_sms_id_fkey";

ALTER TABLE "user_activity" DROP CONSTRAINT IF EXISTS duration_time_user_activity_validation;
ALTER TABLE "send_verification_code" DROP CONSTRAINT IF EXISTS duration_time_send_verification_code_validation;

DROP TABLE IF EXISTS "notification_log";
DROP TABLE IF EXISTS "send_verification_code";
DROP TABLE IF EXISTS "send_email_sms";
DROP TABLE IF EXISTS "user_activity";

DROP TYPE IF EXISTS "activity_status";
DROP TYPE IF EXISTS "notification_log_status";
