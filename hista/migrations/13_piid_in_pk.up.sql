-- modify "notes" table
DROP INDEX "idx_notes_pi_id";
ALTER TABLE "notes" DROP CONSTRAINT "notes_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD PRIMARY KEY ("id", "pi_id");


-- modify "headaches" table
DROP INDEX "idx_headaches_pi_id";
ALTER TABLE "headaches" DROP CONSTRAINT "headaches_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD PRIMARY KEY ("id", "pi_id");


-- modify "pollens" table
CREATE SEQUENCE IF NOT EXISTS "pollens_id_seq" OWNED BY "pollens"."id";-- create sequence for serial column "id"
ALTER TABLE "pollens" ALTER COLUMN "id" SET DEFAULT nextval('"pollens_id_seq"'), ALTER COLUMN "id" DROP DEFAULT;


-- modify "statuses" table
DROP INDEX "idx_statuses_pi_id";
ALTER TABLE "statuses" DROP CONSTRAINT "statuses_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD PRIMARY KEY ("id", "pi_id");


-- modify "condition_events" table
DROP INDEX "idx_condition_events_pi_id";
ALTER TABLE "conditions" DROP CONSTRAINT "fk_condition_events_conditions"; -- enable dropping primary key in "condition_events"
ALTER TABLE "condition_events" DROP CONSTRAINT "condition_events_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD PRIMARY KEY ("id", "pi_id");

-- modify "symptom_categories" table
DROP INDEX "idx_symptom_categories_pi_id";
ALTER TABLE "symptoms" DROP CONSTRAINT "fk_symptom_categories_symptoms"; -- enable dropping primary key in "symptom_categories"
ALTER TABLE "symptom_categories" DROP CONSTRAINT "symptom_categories_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD PRIMARY KEY ("id", "pi_id");

-- modify "symptoms" table
DROP INDEX "idx_symptoms_pi_id"; 
ALTER TABLE "conditions"  DROP CONSTRAINT "fk_conditions_symptom"; -- enable dropping primary key in "symptoms"
ALTER TABLE "symptoms" DROP CONSTRAINT "symptoms_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD COLUMN "symptom_category_pi_id" uuid NULL, ADD PRIMARY KEY ("id", "pi_id"), ADD
 CONSTRAINT "fk_symptom_categories_symptoms" FOREIGN KEY ("symptom_category_id", "symptom_category_pi_id") REFERENCES "symptom_categories" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE NO ACTION;
UPDATE "symptoms" SET "symptom_category_pi_id" = "pi_id";

-- modify "conditions" table
DROP INDEX "idx_conditions_pi_id";
ALTER TABLE "conditions" DROP CONSTRAINT "conditions_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD COLUMN "symptom_pi_id" uuid NULL, ADD COLUMN "condition_event_pi_id" uuid NULL, ADD PRIMARY KEY ("id", "pi_id"), ADD
 CONSTRAINT "fk_condition_events_conditions" FOREIGN KEY ("condition_event_id", "condition_event_pi_id") REFERENCES "condition_events" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD
 CONSTRAINT "fk_conditions_symptom" FOREIGN KEY ("symptom_id", "symptom_pi_id") REFERENCES "symptoms" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE NO ACTION;
UPDATE "conditions" SET "symptom_pi_id" = "pi_id", "condition_event_pi_id" = "pi_id";


-- modify "ingredients" table
DROP INDEX "idx_ingredients_pi_id"; 
ALTER TABLE "foods" DROP CONSTRAINT "fk_foods_ingredient"; -- enable dropping primary key in "ingredients"
ALTER TABLE "ingredients" DROP CONSTRAINT "ingredients_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD PRIMARY KEY ("id", "pi_id");

-- modify "meals" table
DROP INDEX "idx_meals_pi_id"; 
ALTER TABLE "foods" DROP CONSTRAINT "fk_meals_foods"; -- enable dropping primary key in "meals"
ALTER TABLE "meals" DROP CONSTRAINT "meals_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD PRIMARY KEY ("id", "pi_id");

-- modify "foods" table
DROP INDEX "idx_foods_pi_id";
ALTER TABLE "foods" DROP CONSTRAINT "foods_pkey", ALTER COLUMN "pi_id" SET NOT NULL, ADD COLUMN "ingredient_pi_id" uuid NULL, ADD COLUMN "meal_pi_id" uuid NULL, ADD PRIMARY KEY ("id", "pi_id"), ADD
 CONSTRAINT "fk_foods_ingredient" FOREIGN KEY ("ingredient_id", "ingredient_pi_id") REFERENCES "ingredients" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD
 CONSTRAINT "fk_meals_foods" FOREIGN KEY ("meal_id", "meal_pi_id") REFERENCES "meals" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE CASCADE;
UPDATE "foods" SET "ingredient_pi_id" = "pi_id", "meal_pi_id" = "pi_id";


-- drop "tokens" table
DROP TABLE "tokens";
-- drop "users" table
DROP TABLE "users";
