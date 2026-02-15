ALTER TABLE "ingredients" ADD COLUMN "is_archived" boolean;

UPDATE "ingredients" SET "is_archived" = false;

ALTER TABLE "ingredients" ALTER COLUMN "is_archived" SET NOT NULL;