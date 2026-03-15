-- modify "ingredients" table
ALTER TABLE "ingredients" ADD COLUMN "nutrition_protein" bigint NULL, ADD COLUMN "nutrition_carbohydrate" bigint NULL, ADD COLUMN "nutrition_fat" bigint NULL, ADD COLUMN "nutrition_fiber" bigint NULL;
