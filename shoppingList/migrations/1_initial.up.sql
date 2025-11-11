-- create "products" table
CREATE TABLE "products" (
  "id" bigserial NOT NULL,
  "pi_id" uuid NOT NULL,
  "name" text NULL,
  PRIMARY KEY ("id", "pi_id")
);
-- create index "idx_name_piid" to table: "products"
CREATE UNIQUE INDEX "idx_name_piid" ON "products" ("pi_id", "name");
-- create "lists" table
CREATE TABLE "lists" (
  "id" bigserial NOT NULL,
  "pi_id" uuid NOT NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id", "pi_id")
);
-- create "items" table
CREATE TABLE "items" (
  "id" bigserial NOT NULL,
  "pi_id" uuid NOT NULL,
  "product_id" bigint NULL,
  "product_piid" uuid NULL,
  "list_id" bigint NULL,
  "list_piid" uuid NULL,
  "checked" boolean NULL,
  "quantity" smallint NULL,
  "created_at" timestamptz NULL,
  "deleted_at" timestamptz NULL,
  PRIMARY KEY ("id", "pi_id"),
  CONSTRAINT "fk_items_product" FOREIGN KEY ("product_id", "product_piid") REFERENCES "products" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_lists_items" FOREIGN KEY ("list_id", "list_piid") REFERENCES "lists" ("id", "pi_id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_product_list" to table: "items"
CREATE UNIQUE INDEX "idx_product_list" ON "items" ("product_id", "list_id");
