-- drop index "idx_ingredients_name" from table: "ingredients"
DROP INDEX "idx_ingredients_name";
-- create index "idx_name_piid" to table: "ingredients"
CREATE UNIQUE INDEX "idx_name_piid" ON "ingredients" ("pi_id", "name");

-- create "templates" table
CREATE TABLE "templates" (
  "id" bigserial NOT NULL,
  "pi_id" uuid NOT NULL,
  "name" text NULL,
  PRIMARY KEY ("id", "pi_id")
);
-- create index "idx_template_name_piid" to table: "templates"
CREATE UNIQUE INDEX "idx_template_name_piid" ON "templates" ("pi_id", "name");
-- create "template_items" table
CREATE TABLE "template_items" (
  "id" bigserial NOT NULL,
  "pi_id" uuid NOT NULL,
  "condition" text NULL,
  "template_id" bigint NULL,
  "template_pi_id" uuid NULL,
  "ingredient_id" bigint NULL,
  "ingredient_pi_id" uuid NULL,
  PRIMARY KEY ("id", "pi_id"),
  CONSTRAINT "fk_ingredients_items" FOREIGN KEY ("ingredient_id", "ingredient_pi_id") REFERENCES "ingredients" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE CASCADE,
  CONSTRAINT "fk_templates_items" FOREIGN KEY ("template_id", "template_pi_id") REFERENCES "templates" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE CASCADE
);
