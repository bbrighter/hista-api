-- reverse: modify "statuses" table
ALTER TABLE "statuses" DROP COLUMN "morning_sleep", DROP COLUMN "evening_fitness", DROP COLUMN "morning_fitness";

-- reverse: modify "morning_statuses" table
-- reverse drop "evening_statuses" table
CREATE TABLE "evening_statuses" (
  "id" bigserial NOT NULL,
  "status_id" bigint NULL,
  "fitness" smallint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_evening_statuses_status_id" UNIQUE ("status_id"),
  CONSTRAINT "fk_statuses_evening" FOREIGN KEY ("status_id") REFERENCES "statuses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
-- reverse drop "morning_statuses" table
CREATE TABLE "morning_statuses" (
  "id" bigserial NOT NULL,
  "status_id" bigint NULL,
  "fitness" smallint NULL,
  "sleep" smallint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "uni_morning_statuses_status_id" UNIQUE ("status_id"),
  CONSTRAINT "fk_statuses_morning" FOREIGN KEY ("status_id") REFERENCES "statuses" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);

