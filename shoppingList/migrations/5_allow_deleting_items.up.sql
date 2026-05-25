-- modify "items" table
DELETE FROM "items" WHERE "deleted_at" IS NOT NULL;

ALTER TABLE "items" DROP COLUMN "deleted_at";
