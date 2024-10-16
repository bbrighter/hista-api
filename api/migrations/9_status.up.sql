-- create "statuses" table
CREATE TABLE "statuses" (
  "id" bigserial NOT NULL,
  "date" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "evening_statuses" table
CREATE TABLE "evening_statuses" (
  "id" bigserial NOT NULL,
  "status_id" bigint NULL,
  "fitness" smallint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_statuses_evening" FOREIGN KEY ("status_id") REFERENCES "statuses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- create "morning_statuses" table
CREATE TABLE "morning_statuses" (
  "id" bigserial NOT NULL,
  "status_id" bigint NULL,
  "fitness" smallint NULL,
  "sleep" smallint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_statuses_morning" FOREIGN KEY ("status_id") REFERENCES "statuses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
