-- create "product_instances" table
CREATE TABLE "product_instances" (
  "id" uuid NOT NULL,
  "name" text NULL,
  "product_id" text NULL,
  PRIMARY KEY ("id")
);
