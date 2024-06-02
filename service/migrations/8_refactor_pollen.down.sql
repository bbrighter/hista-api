-- reverse: drop "pollens" table
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
-- reverse: create "pollen_v2" table
DROP TABLE "pollen_v2";
-- reverse: create "pollen_events" table
DROP TABLE "pollen_events";
