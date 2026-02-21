-- create "medicines" table
CREATE TABLE "medicines" (
  "id" bigserial NOT NULL,
  "pi_id" uuid NOT NULL,
  "name" text NULL,
  "is_archived" boolean NOT NULL,
  PRIMARY KEY ("id", "pi_id")
);
-- create index "idx_medicines_name" to table: "medicines"
CREATE UNIQUE INDEX "idx_medicines_name" ON "medicines" ("name");
-- create "intakes" table
CREATE TABLE "intakes" (
  "id" bigserial NOT NULL,
  "pi_id" uuid NOT NULL,
  "date" timestamptz NOT NULL,
  "medicine_id" bigint NULL,
  "medicine_pi_id" uuid NULL,
  PRIMARY KEY ("id", "pi_id"),
  CONSTRAINT "fk_medicines_intakes" FOREIGN KEY ("medicine_id", "medicine_pi_id") REFERENCES "medicines" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
