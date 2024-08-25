-- modify "meals" table
ALTER TABLE "meals" ADD COLUMN "freshness" smallint NULL, ADD COLUMN "stress_level" smallint NULL, ADD COLUMN "is_alone" boolean NULL;

UPDATE "meals" 
SET "freshness" = 0, "stress_level" = 0, "is_alone" = true;
