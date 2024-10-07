-- create "statuses" table
CREATE TABLE "statuses" (
  "id" bigserial NOT NULL,
  "date" timestamptz NULL,
  "time_of_day" text NULL,
  "fitness" smallint NULL,
  "sleep" smallint NULL,
  PRIMARY KEY ("id")
);
