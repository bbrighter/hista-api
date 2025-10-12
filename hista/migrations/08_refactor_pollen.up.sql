-- create "pollen_events" table
CREATE TABLE "pollen_events" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);

-- rename old pollens table and drop later
ALTER TABLE "pollens" 
RENAME TO "pollen_v2";

-- create "pollen_v2" table
CREATE TABLE "pollens" (
  "id" bigserial NOT NULL,
  "pollen_event_id" bigint NULL,
  "type" text NULL,
  "intensity" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_pollen_events_pollens" FOREIGN KEY ("pollen_event_id") REFERENCES "pollen_events" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);

-- migrate existing events
INSERT INTO "pollen_events" ("created_at")
  SELECT "created_at" FROM "pollen_v2";

-- migrate existing pollen based on creation_date
INSERT INTO "pollens" ("pollen_event_id", "type", "intensity")
  SELECT "pollen_events"."id", 'Roggen', "pollen_v2"."roggen" 
  FROM "pollen_v2" JOIN "pollen_events" ON "pollen_v2"."created_at" = "pollen_events"."created_at";

INSERT INTO "pollens" ("pollen_event_id", "type", "intensity")
  SELECT "pollen_events"."id", 'Ambrosia', "pollen_v2"."ambrosia" 
  FROM "pollen_v2" JOIN "pollen_events" ON "pollen_v2"."created_at" = "pollen_events"."created_at";

INSERT INTO "pollens" ("pollen_event_id", "type", "intensity")
  SELECT "pollen_events"."id", 'Erle', "pollen_v2"."erle" 
  FROM "pollen_v2" JOIN "pollen_events" ON "pollen_v2"."created_at" = "pollen_events"."created_at";

INSERT INTO "pollens" ("pollen_event_id", "type", "intensity")
  SELECT "pollen_events"."id", 'Beifuss', "pollen_v2"."beifuss" 
  FROM "pollen_v2" JOIN "pollen_events" ON "pollen_v2"."created_at" = "pollen_events"."created_at";

INSERT INTO "pollens" ("pollen_event_id", "type", "intensity")
  SELECT "pollen_events"."id", 'Birke', "pollen_v2"."birke" 
  FROM "pollen_v2" JOIN "pollen_events" ON "pollen_v2"."created_at" = "pollen_events"."created_at";

INSERT INTO "pollens" ("pollen_event_id", "type", "intensity")
  SELECT "pollen_events"."id", 'Gräser', "pollen_v2"."graeser" 
  FROM "pollen_v2" JOIN "pollen_events" ON "pollen_v2"."created_at" = "pollen_events"."created_at";

INSERT INTO "pollens" ("pollen_event_id", "type", "intensity")
  SELECT "pollen_events"."id", 'Hasel', "pollen_v2"."hasel" 
  FROM "pollen_v2" JOIN "pollen_events" ON "pollen_v2"."created_at" = "pollen_events"."created_at";

INSERT INTO "pollens" ("pollen_event_id", "type", "intensity")
  SELECT "pollen_events"."id", 'Esche', "pollen_v2"."esche" 
  FROM "pollen_v2" JOIN "pollen_events" ON "pollen_v2"."created_at" = "pollen_events"."created_at";


-- drop old pollens table
DROP TABLE "pollen_v2";

