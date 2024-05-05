-- create "symptom_categories" table
CREATE TABLE "symptom_categories" (
  "id" bigserial NOT NULL,
  "name" text NULL,
  PRIMARY KEY ("id")
);
-- create "symptoms" table
CREATE TABLE "symptoms" (
  "id" bigserial NOT NULL,
  "name" text NULL,
  "symptom_category_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_symptom_categories_symptoms" FOREIGN KEY ("symptom_category_id") REFERENCES "symptom_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create "states" table
CREATE TABLE "states" (
  "id" bigserial NOT NULL,
  "date" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "conditions" table
CREATE TABLE "conditions" (
  "id" bigserial NOT NULL,
  "symptom_id" bigint NULL,
  "severity" smallint NULL,
  "state_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_conditions_symptom" FOREIGN KEY ("symptom_id") REFERENCES "symptoms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_states_conditions" FOREIGN KEY ("state_id") REFERENCES "states" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
