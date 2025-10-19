-- create "products" table
CREATE TABLE "products" (
  "id" bigserial NOT NULL,
  "pi_id" uuid NOT NULL,
  "name" text NULL,
  PRIMARY KEY ("id", "pi_id")
);
-- create index "idx_name_piid" to table: "products"
CREATE UNIQUE INDEX "idx_name_piid" ON "products" ("pi_id", "name");
-- create "items" table
CREATE TABLE "items" (
  "id" bigserial NOT NULL,
  "pi_id" text NOT NULL,
  "product_id" bigint NULL,
  "product_piid" uuid NULL,
  "checked" boolean NULL,
  "quantity" smallint NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id", "pi_id"),
  CONSTRAINT "fk_items_product" FOREIGN KEY ("product_id", "product_piid") REFERENCES "products" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
