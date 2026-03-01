ALTER TABLE "medicines" ADD COLUMN "sort_order" bigint NULL;

WITH reordered as (
    SELECT 
        "id",
        "pi_id",
        ROW_NUMBER() OVER (PARTITION BY "pi_id" ORDER BY "name") * 100 as "sort_order"
    FROM "medicines"
)
UPDATE "medicines" m 
SET "sort_order" = r."sort_order"
FROM reordered r
WHERE m."id" = r."id" AND m."pi_id" = r."pi_id";
 
ALTER TABLE "medicines" ALTER COLUMN "sort_order" SET NOT NULL;

