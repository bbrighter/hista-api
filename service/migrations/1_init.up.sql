-- create "ingredients" table
CREATE TABLE "ingredients" (
  "id" bigserial NOT NULL,
  "name" text NULL,
  PRIMARY KEY ("id")
);
-- create "meals" table
CREATE TABLE "meals" (
  "id" bigserial NOT NULL,
  "date" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- create "foods" table
CREATE TABLE "foods" (
  "id" bigserial NOT NULL,
  "ingredient_id" bigint NULL,
  "condition" text NULL,
  "meal_id" bigint NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_foods_ingredient" FOREIGN KEY ("ingredient_id") REFERENCES "ingredients" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_meals_foods" FOREIGN KEY ("meal_id") REFERENCES "meals" ("id") ON UPDATE NO ACTION ON DELETE CASCADE
);
