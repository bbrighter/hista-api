-- modify "lists" table
ALTER TABLE "lists" ADD COLUMN "singleton" boolean NOT NULL DEFAULT true;
-- create index "idx_lists_singleton" to table: "lists"
CREATE UNIQUE INDEX "idx_lists_singleton" ON "lists" ("pi_id", "singleton") WHERE (deleted_at IS NULL);
-- drop "moments" table
DROP TABLE "moments";
