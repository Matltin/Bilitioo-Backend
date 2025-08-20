CREATE TABLE "verify_emails" (
  "id" bigserial PRIMARY KEY,
  "user_id" bigint NOT NULL,
  "email" varchar NOT NULL,
  "secret_code" varchar NOT NULL,
  "is_used" bool NOT NULL DEFAULT false,
  "created_at" timestamptz NOT NULL DEFAULT (now()),
  "expired_at" timestamptz NOT NULL DEFAULT (now() + interval '15 minutes')
);

ALTER TABLE "verify_emails" ADD FOREIGN KEY ("user_id") REFERENCES "user" ("id");

-- Create random admin users
INSERT INTO "user" (
  "email", 
  "phone_number", 
  "hashed_password", 
  "role", 
  "status", 
  "phone_verified", 
  "email_verified"
) VALUES 
-- Admin users
('admin@transport.ir', '09123456789', 'test_password', 'ADMIN', 'ACTIVE', true, true),
('manager@transport.ir', '09123456788', 'test_password', 'ADMIN', 'ACTIVE', true, true),
('support@transport.ir', '09123456787', 'test_password', 'ADMIN', 'ACTIVE', true, true);