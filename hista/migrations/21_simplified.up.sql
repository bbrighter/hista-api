---- Symptoms ----
-- Drop FK
ALTER TABLE "conditions" DROP CONSTRAINT "fk_condition_events_conditions", DROP CONSTRAINT "fk_conditions_symptom";
ALTER TABLE "symptoms" DROP CONSTRAINT "fk_symptom_categories_symptoms";
-- Drop PK
ALTER TABLE "condition_events" DROP CONSTRAINT "condition_events_pkey";
ALTER TABLE "conditions" DROP CONSTRAINT "conditions_pkey";
ALTER TABLE "symptoms" DROP CONSTRAINT "symptoms_pkey";
ALTER TABLE "symptom_categories" DROP CONSTRAINT "symptom_categories_pkey";
-- Drop PIID
ALTER TABLE "conditions" DROP COLUMN "pi_id", DROP COLUMN "symptom_pi_id", DROP COLUMN "condition_event_pi_id";
ALTER TABLE "symptoms" DROP COLUMN "pi_id", DROP COLUMN "symptom_category_pi_id";
-- Add PK
ALTER TABLE "condition_events" ADD PRIMARY KEY ("id");
ALTER TABLE "conditions" ADD PRIMARY KEY ("id");
ALTER TABLE "symptoms" ADD PRIMARY KEY ("id");
ALTER TABLE "symptom_categories" ADD PRIMARY KEY ("id");
-- Add indices
CREATE INDEX "idx_condition_events_pi_id" ON "condition_events" ("pi_id");
CREATE INDEX "idx_symptom_categories_pi_id" ON "symptom_categories" ("pi_id");
-- Add FK
ALTER TABLE "conditions"
    ADD CONSTRAINT "fk_condition_events_conditions"
    FOREIGN KEY ("condition_event_id") REFERENCES "condition_events" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "conditions"
    ADD CONSTRAINT "fk_conditions_symptom"
    FOREIGN KEY ("symptom_id") REFERENCES "symptoms" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "symptoms"
    ADD CONSTRAINT "fk_symptom_categories_symptoms"
    FOREIGN KEY ("symptom_category_id") REFERENCES "symptom_categories" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
---------------

---- Meals ----
-- Drop FKs
ALTER TABLE "foods" DROP CONSTRAINT "fk_foods_ingredient", DROP CONSTRAINT "fk_meals_foods";
ALTER TABLE "template_items" DROP CONSTRAINT "fk_ingredients_items", DROP CONSTRAINT "fk_templates_items";
-- Drop PKs
ALTER TABLE "meals" DROP CONSTRAINT "meals_pkey";
ALTER TABLE "ingredients" DROP CONSTRAINT "ingredients_pkey";
ALTER TABLE "foods" DROP CONSTRAINT "foods_pkey";
ALTER TABLE "templates" DROP CONSTRAINT "templates_pkey";
ALTER TABLE "template_items" DROP CONSTRAINT "template_items_pkey";
-- Drop PIID
ALTER TABLE "foods" DROP COLUMN "pi_id", DROP COLUMN "ingredient_pi_id", DROP COLUMN "meal_pi_id";
ALTER TABLE "template_items" DROP COLUMN "template_pi_id", DROP COLUMN "ingredient_pi_id", DROP COLUMN "pi_id";
-- Add PK
ALTER TABLE "meals" ADD PRIMARY KEY ("id");
ALTER TABLE "ingredients" ADD PRIMARY KEY ("id");
ALTER TABLE "foods" ADD PRIMARY KEY ("id");
ALTER TABLE "templates" ADD PRIMARY KEY ("id");
ALTER TABLE "template_items" ADD PRIMARY KEY ("id");
-- Add index
CREATE INDEX "idx_templates_pi_id" ON "templates" ("pi_id");
CREATE INDEX "idx_meals_pi_id" ON "meals" ("pi_id");
CREATE INDEX "idx_ingredients_pi_id" ON "ingredients" ("pi_id");
-- Add FK
ALTER TABLE "foods" 
    ADD CONSTRAINT "fk_foods_ingredient" 
    FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
ALTER TABLE "foods" 
    ADD CONSTRAINT "fk_meals_foods" 
    FOREIGN KEY ("meal_id") REFERENCES "meals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "template_items" 
    ADD CONSTRAINT "fk_ingredients_items" 
    FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
ALTER TABLE "template_items" 
    ADD CONSTRAINT "fk_templates_items" 
    FOREIGN KEY ("template_id") REFERENCES "templates" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
-------------------

---- Medicines ----
-- Drop FK
ALTER TABLE "intakes" DROP CONSTRAINT "fk_medicines_intakes";
-- Drop PK
ALTER TABLE "medicines" DROP CONSTRAINT "medicines_pkey";
ALTER TABLE "intakes" DROP CONSTRAINT "intakes_pkey";
-- Drop PIID
ALTER TABLE "intakes" DROP COLUMN "pi_id", DROP COLUMN "medicine_pi_id";
-- Add PK
ALTER TABLE "medicines" ADD PRIMARY KEY ("id");
ALTER TABLE "intakes" ADD PRIMARY KEY ("id");
-- Add index
CREATE INDEX "idx_medicines_pi_id" ON "medicines" ("pi_id");
-- Add FK
ALTER TABLE "intakes"
    ADD CONSTRAINT "fk_medicines_intakes"
    FOREIGN KEY ("medicine_id") REFERENCES "medicines" ("id") ON UPDATE NO ACTION ON DELETE CASCADE;
---------------

---- Notes ----
-- Drop PK
ALTER TABLE "notes" DROP CONSTRAINT "notes_pkey";
-- Add PK
ALTER TABLE "notes" ADD PRIMARY KEY ("id");
-- Add index
CREATE INDEX "idx_notes_pi_id" ON "notes" ("pi_id");
---------------


---- Headaches ----
-- Drop PK
ALTER TABLE "headaches" DROP CONSTRAINT "headaches_pkey";
-- Add PK
ALTER TABLE "headaches" ADD PRIMARY KEY ("id");
-- Add index
CREATE INDEX "idx_headaches_pi_id" ON "headaches" ("pi_id");


---- Statuses ----
-- Drop PK
ALTER TABLE "statuses" DROP CONSTRAINT "statuses_pkey";
-- Add PK
ALTER TABLE "statuses" ADD PRIMARY KEY ("id");
-- Add index
CREATE INDEX "idx_statuses_pi_id" ON "statuses" ("pi_id");