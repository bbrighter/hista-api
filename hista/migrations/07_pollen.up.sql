-- create "pollens" table
CREATE TABLE "pollens" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "roggen" bigint NULL,
  "ambrosia" bigint NULL,
  "erle" bigint NULL,
  "beifuss" bigint NULL,
  "birke" bigint NULL,
  "graeser" bigint NULL,
  "hasel" bigint NULL,
  "esche" bigint NULL,
  PRIMARY KEY ("id")
);
