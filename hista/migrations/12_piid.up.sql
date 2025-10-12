-- modify "notes" table
ALTER TABLE "notes" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_notes_pi_id" to table: "notes"
CREATE INDEX "idx_notes_pi_id" ON "notes" ("pi_id");
-- modify "statuses" table
ALTER TABLE "statuses" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_statuses_pi_id" to table: "statuses"
CREATE INDEX "idx_statuses_pi_id" ON "statuses" ("pi_id");
-- modify "foods" table
ALTER TABLE "foods" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_foods_pi_id" to table: "foods"
CREATE INDEX "idx_foods_pi_id" ON "foods" ("pi_id");
-- modify "headaches" table
ALTER TABLE "headaches" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_headaches_pi_id" to table: "headaches"
CREATE INDEX "idx_headaches_pi_id" ON "headaches" ("pi_id");
-- modify "ingredients" table
ALTER TABLE "ingredients" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_ingredients_pi_id" to table: "ingredients"
CREATE INDEX "idx_ingredients_pi_id" ON "ingredients" ("pi_id");
-- modify "meals" table
ALTER TABLE "meals" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_meals_pi_id" to table: "meals"
CREATE INDEX "idx_meals_pi_id" ON "meals" ("pi_id");
-- modify "conditions" table
ALTER TABLE "conditions" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_conditions_pi_id" to table: "conditions"
CREATE INDEX "idx_conditions_pi_id" ON "conditions" ("pi_id");
-- create sequence for serial column "id"
-- modify "condition_events" table
ALTER TABLE "condition_events" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_condition_events_pi_id" to table: "condition_events"
CREATE INDEX "idx_condition_events_pi_id" ON "condition_events" ("pi_id");
-- modify "symptom_categories" table
ALTER TABLE "symptom_categories" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_symptom_categories_pi_id" to table: "symptom_categories"
CREATE INDEX "idx_symptom_categories_pi_id" ON "symptom_categories" ("pi_id");
-- modify "symptoms" table
ALTER TABLE "symptoms" ADD COLUMN "pi_id" uuid NULL;
-- create index "idx_symptoms_pi_id" to table: "symptoms"
CREATE INDEX "idx_symptoms_pi_id" ON "symptoms" ("pi_id");
