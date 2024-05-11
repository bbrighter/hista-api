-- create "users" table
CREATE TABLE "users" (
  "id" bigserial NOT NULL,
  "name" text NULL,
  "password_hash_hex" text NULL,
  PRIMARY KEY ("id")
);
-- create "tokens" table
CREATE TABLE "tokens" (
  "id" bigserial NOT NULL,
  "bearer" text NULL,
  "expires" timestamptz NULL,
  "user_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_users_tokens" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
