CREATE TYPE "role" AS ENUM (
  'ADMIN',
  'USER'
);

CREATE TYPE "user_status" AS ENUM (
  'ACTIVE',
  'NON-ACTIVE'
);

CREATE TABLE "user" (
  "id" bigserial PRIMARY KEY,
  "email" varchar NOT NULL,
  "phone_number" varchar(11) NOT NULL,
  "hashed_password" varchar NOT NULL,
  "password_change_at" timestamptz NOT NULL DEFAULT (now()),
  "role" role NOT NULL DEFAULT 'USER',
  "status" user_status NOT NULL DEFAULT 'ACTIVE',
  "phone_verified" bool NOT NULL DEFAULT false,
  "email_verified" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  UNIQUE("email", "phone_number")
);

CREATE TABLE "profile" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint UNIQUE NOT NULL,
  "pic_dir" varchar  NOT NULL,
  "first_name" varchar NOT NULL,
  "last_name" varchar  NOT NULL,
  "city_id" bigint NOT NULL,
  "national_code" varchar NOT NULL
);


CREATE TABLE "city" (
  "id" bigserial PRIMARY KEY,
  "province" varchar NOT NULL,
  "county" varchar NOT NULL
);


ALTER TABLE "user"
ADD CONSTRAINT email_or_phone_required
CHECK (
    (email IS NOT NULL AND email <> '') OR 
    (phone_number IS NOT NULL AND phone_number <> '')
);

CREATE INDEX ON "user" ("email");
CREATE INDEX ON "user" ("phone_number");
CREATE INDEX ON "profile" ("user_id");
CREATE INDEX ON "city" ("province");
CREATE INDEX ON "city" ("county");
CREATE INDEX ON "city" ("province", "county");

ALTER TABLE "profile" 
ADD CONSTRAINT profile_user_id_fkey
FOREIGN KEY ("user_id") 
REFERENCES "user" ("id")
ON DELETE CASCADE;


ALTER TABLE "profile"
ADD CONSTRAINT profile_city_id_fkey
FOREIGN KEY ("city_id")
REFERENCES "city" ("id")
ON DELETE SET NULL;

CREATE TYPE "check_reservation_ticket_status" AS ENUM (
    'RESERVED',
    'NOT_RESERVED'
);

CREATE TYPE "vehicle_type" AS ENUM (
  'BUS',
  'TRAIN',
  'AIRPLANE'
);

CREATE TYPE "flight_class" AS ENUM (
  'ECONOMY',
  'PREMIUM-ECONOMY',
  'BUSINESS',
  'FIRST'
);


CREATE TABLE "ticket" (
  "id" bigserial PRIMARY KEY,
  "vehicle_id" bigint NOT NULL,
  "seat_id" bigint NOT NULL,
  "vehicle_type" vehicle_type NOT NULL,
  "route_id" bigint NOT NULL,
  "amount" bigint NOT NULL,
  "departure_time" timestamptz NOT NULL,
  "arrival_time" timestamptz NOT NULL,
  "count_stand" int NOT NULL DEFAULT 0,
  "status" check_reservation_ticket_status NOT NULL DEFAULT 'NOT_RESERVED',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "vehicle" (
  "id" bigserial PRIMARY KEY,
  "company_id" bigint NOT NULL,
  "capacity" int NOT NULL,
  "vehicle_type" vehicle_type NOT NULL,
  "feature" json NOT NULL
);

CREATE TABLE "company" (
  "id" bigserial PRIMARY KEY,
  "name" text NOT NULL,
  "address" text NOT NULL
);

CREATE TABLE "bus" (
  "vehicle_id" bigint NOT NULL,
  "VIP" boolean NOT NULL DEFAULT false,
  "bed_chair" boolean NOT NULL DEFAULT false
);

CREATE TABLE "train" (
  "vehicle_id" bigint NOT NULL,
  "rank" int NOT NULL DEFAULT 3,
  "have_compartment" boolean NOT NULL DEFAULT false
);

CREATE TABLE "airplane" (
  "vehicle_id" bigint NOT NULL,
  "flight_class" flight_class NOT NULL,
  "name" varchar NOT NULL
);

CREATE INDEX ON "ticket" ("route_id");

CREATE INDEX ON "ticket" ("departure_time");

CREATE INDEX ON "ticket" ("route_id", "departure_time", "vehicle_id");

CREATE INDEX ON "vehicle" ("company_id");

ALTER TABLE "ticket" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicle" ("id");

ALTER TABLE "vehicle" ADD FOREIGN KEY ("company_id") REFERENCES "company" ("id");

ALTER TABLE "bus" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicle" ("id");

ALTER TABLE "train" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicle" ("id");

ALTER TABLE "airplane" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicle" ("id");

ALTER TABLE "ticket" ADD CONSTRAINT amount_validation CHECK (amount > 0);

ALTER TABLE "ticket" ADD CONSTRAINT time_validation CHECK (arrival_time > departure_time);

ALTER TABLE "ticket" ADD CONSTRAINT count_stand_validation CHECK (count_stand >= 0);

ALTER TABLE "vehicle" ADD CONSTRAINT capacity_validation CHECK (capacity > 0);

ALTER TABLE "train" ADD CONSTRAINT rank_validation CHECK (rank BETWEEN 3 AND 5);

CREATE TABLE "route" (
  "id" bigserial PRIMARY KEY,
  "origin_city_id" bigint NOT NULL,
  "destination_city_id" bigint NOT NULL,
  "origin_terminal_id" bigint,
  "destination_terminal_id" bigint,
  "distance" int NOT NULL
);

CREATE TABLE "terminal" (
  "id" bigserial PRIMARY KEY,
  "name" text NOT NULL,
  "address" text NOT NULL
);

CREATE INDEX ON "route" ("origin_city_id");
CREATE INDEX ON "route" ("destination_city_id");
CREATE INDEX ON "route" ("origin_terminal_id");
CREATE INDEX ON "route" ("origin_city_id", "destination_city_id");
CREATE INDEX ON "route" ("destination_terminal_id");
CREATE INDEX ON "route" ("origin_terminal_id", "destination_terminal_id");

ALTER TABLE "ticket"
ADD CONSTRAINT ticket_route_id_fkey
FOREIGN KEY ("route_id") REFERENCES "route"("id") ON DELETE CASCADE;

ALTER TABLE "route"
ADD CONSTRAINT route_origin_city_id_fkey
FOREIGN KEY ("origin_city_id") REFERENCES "city"("id") ON DELETE CASCADE;

ALTER TABLE "route"
ADD CONSTRAINT route_destination_city_id_fkey
FOREIGN KEY ("destination_city_id") REFERENCES "city"("id") ON DELETE CASCADE;

ALTER TABLE "route"
ADD CONSTRAINT route_origin_terminal_id_fkey
FOREIGN KEY ("origin_terminal_id") REFERENCES "terminal"("id") ON DELETE CASCADE;

ALTER TABLE "route"
ADD CONSTRAINT route_destination_terminal_id_fkey
FOREIGN KEY ("destination_terminal_id") REFERENCES "terminal"("id") ON DELETE CASCADE;

ALTER TABLE "route"
ADD CONSTRAINT distance_validation
CHECK (distance > 0);

CREATE TABLE "seat" (
  "id" bigserial PRIMARY KEY,
  "vehicle_id" bigint NOT NULL,
  "vehicle_type" vehicle_type NOT NULL,
  "seat_number" int NOT NULL,
  "is_available" boolean NOT NULL DEFAULT true
);

CREATE TABLE "bus_seat" (
  "seat_id" bigint NOT NULL
);

CREATE TABLE "train_seat" (
  "seat_id" bigint NOT NULL,
  "salon" int NOT NULL,
  "have_compartment" boolean NOT NULL DEFAULT false,
  "cuope_number" int
);

CREATE TABLE "airplane_seat" (
  "seat_id" bigint NOT NULL
);

CREATE INDEX ON "seat" ("vehicle_id");

ALTER TABLE "ticket" ADD FOREIGN KEY ("seat_id") REFERENCES "seat" ("id");

ALTER TABLE "seat" ADD FOREIGN KEY ("vehicle_id") REFERENCES "vehicle" ("id");

ALTER TABLE "bus_seat" ADD FOREIGN KEY ("seat_id") REFERENCES "seat" ("id");

ALTER TABLE "train_seat" ADD FOREIGN KEY ("seat_id") REFERENCES "seat" ("id");

ALTER TABLE "airplane_seat" ADD FOREIGN KEY ("seat_id") REFERENCES "seat" ("id");

ALTER TABLE "seat" ADD CONSTRAINT seat_number_validation CHECK (seat_number > 0);

ALTER TABLE "train_seat" ADD CONSTRAINT salon_validation CHECK (salon > 0);

ALTER TABLE "train_seat" ADD CONSTRAINT coupe_number_validation CHECK (cuope_number > 0);

CREATE TYPE "request_type" AS ENUM (
  'PAYMENT-ISSUE',
  'TRAVEL-DELAY',
  'UNEXPECTED-RESERVED',
  'ETC.'
);

CREATE TYPE "ticket_status" AS ENUM (
  'RESERVED',
  'RESERVING',
  'CANCELED',
  'CANCELED-BY-TIME'
);

CREATE TYPE "payment_type" AS ENUM (
  'CASH',
  'CREDIT_CARD',
  'WALLET',
  'BANK_TRANSFER',
  'CRYPTO'
);

CREATE TYPE "payment_status" AS ENUM (
  'PENDING',
  'COMPLETED',
  'FAILED',
  'REFUNDED'
);

CREATE TABLE "reservation" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "ticket_id" bigint NOT NULL,
  "payment_id" bigint,
  "status" ticket_status NOT NULL DEFAULT 'RESERVING',
  "duration_time" interval NOT NULL DEFAULT '10 minutes',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "payment" (
  "id" bigserial PRIMARY KEY,
  "from_account" varchar NOT NULL,
  "to_account" varchar NOT NULL,
  "amount" bigint NOT NULL,
  "type" payment_type NOT NULL,
  "status" payment_status NOT NULL DEFAULT 'PENDING',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "change_reservation" (
  "id" bigserial PRIMARY KEY,
  "reservation_id" bigint NOT NULL,
  "admin_id" bigint,
  "user_id" bigint NOT NULL,
  "from_status" ticket_status NOT NULL,
  "to_status" ticket_status NOT NULL
);

CREATE TABLE "report" (
  "id" bigserial PRIMARY KEY,
  "reservation_id" bigint NOT NULL,
  "user_id" bigint NOT NULL,
  "admin_id" bigint NOT NULL,
  "request_type" request_type NOT NULL DEFAULT 'ETC.',
  "request_text" text NOT NULL,
  "response_text" text NOT NULL
);

CREATE INDEX ON "reservation" ("user_id");

CREATE INDEX ON "reservation" ("ticket_id");

CREATE INDEX ON "reservation" ("user_id", "ticket_id");

CREATE INDEX ON "report" ("user_id");

CREATE INDEX ON "report" ("admin_id");

CREATE INDEX ON "report" ("user_id", "admin_id");

ALTER TABLE "reservation"
ADD CONSTRAINT reservation_user_id_fkey
FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE;

ALTER TABLE "reservation"
ADD CONSTRAINT reservation_ticket_id_fkey
FOREIGN KEY ("ticket_id") REFERENCES "ticket"("id") ON DELETE CASCADE;

ALTER TABLE "reservation" ADD FOREIGN KEY ("payment_id") REFERENCES "payment" ("id");

ALTER TABLE "change_reservation"
ADD CONSTRAINT change_reservation_reservation_id_fkey
FOREIGN KEY ("reservation_id") REFERENCES "reservation"("id") ON DELETE CASCADE;

ALTER TABLE "change_reservation"
ADD CONSTRAINT change_reservation_admin_id_fkey
FOREIGN KEY ("admin_id") REFERENCES "user"("id") ON DELETE CASCADE;

ALTER TABLE "change_reservation"
ADD CONSTRAINT change_reservation_user_id_fkey
FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE;

ALTER TABLE "payment" ADD CONSTRAINT amount_payment_validation CHECK (amount > 0);

ALTER TABLE "report"
ADD CONSTRAINT report_reservation_id_fkey
FOREIGN KEY ("reservation_id") REFERENCES "reservation"("id") ON DELETE CASCADE;

ALTER TABLE "report"
ADD CONSTRAINT report_admin_id_fkey
FOREIGN KEY ("admin_id") REFERENCES "user"("id") ON DELETE CASCADE;

ALTER TABLE "report"
ADD CONSTRAINT report_user_id_fkey
FOREIGN KEY ("user_id") REFERENCES "user"("id") ON DELETE CASCADE;

CREATE TABLE "penalty" (
  "id" bigserial PRIMARY KEY,
  "vehicle_id" bigint NOT NULL,
  "penalty_text" text NOT NULL
);

ALTER TABLE "penalty" 
ADD FOREIGN KEY ("vehicle_id") 
REFERENCES "vehicle" ("id");

CREATE TYPE "activity_status" AS ENUM (
  'PENDING',
  'REMINDER-SENT',
  'PURCHASED'
);

CREATE TYPE "notification_log_status" AS ENUM (
  'SENT',
  'FAILED'
);

CREATE TABLE "user_activity" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "route_id" bigint NOT NULL,
  "vehicle_type" vehicle_type NOT NULL,
  "status" activity_status NOT NULL DEFAULT 'PENDING',
  "duration_time" interval NOT NULL DEFAULT '10 minutes',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

CREATE TABLE "notification_log" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "send_email_sms_id" bigint NOT NULL,
  "message_text" text NOT NULL,
  "user_activity_id" bigint NOT NULL
);

CREATE TABLE "send_email_sms" (
  "id" bigserial PRIMARY KEY,
  "message_type" varchar NOT NULL,
  "sent_at" timestamptz NOT NULL DEFAULT (now()),
  "status" notification_log_status NOT NULL DEFAULT 'SENT'
);

CREATE TABLE "send_verification_code" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "send_email_sms_id" bigint NOT NULL,
  "token" varchar NOT NULL,
  "duration_time" interval NOT NULL DEFAULT '10 minutes',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);

ALTER TABLE "user_activity" 
ADD FOREIGN KEY ("user_id") 
REFERENCES "user" ("id");

ALTER TABLE "user_activity" 
ADD FOREIGN KEY ("route_id") 
REFERENCES "route" ("id");

ALTER TABLE "notification_log" 
ADD FOREIGN KEY ("user_id") 
REFERENCES "user" ("id");

ALTER TABLE "notification_log" 
ADD FOREIGN KEY ("send_email_sms_id") 
REFERENCES "send_email_sms" ("id");

ALTER TABLE "notification_log" 
ADD FOREIGN KEY ("user_activity_id") 
REFERENCES "user_activity" ("id");

ALTER TABLE "send_verification_code" 
ADD FOREIGN KEY ("user_id") 
REFERENCES "user" ("id");

ALTER TABLE "send_verification_code" 
ADD FOREIGN KEY ("send_email_sms_id") 
REFERENCES "send_email_sms" ("id");

ALTER TABLE "user_activity"
ADD CONSTRAINT duration_time_user_activity_validation
CHECK (duration_time > interval '0 seconds');

ALTER TABLE "send_verification_code"
ADD CONSTRAINT duration_time_send_verification_code_validation
CHECK (duration_time > interval '0 seconds');