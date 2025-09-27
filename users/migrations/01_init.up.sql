-- create "users" table
CREATE TABLE "users" (
  "id" uuid NOT NULL,
  "name" text NULL,
  "password" text NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_users_name" to table: "users"
CREATE UNIQUE INDEX "idx_users_name" ON "users" ("name");
-- create "user_product_instances" table
CREATE TABLE "user_product_instances" (
  "id" bigserial NOT NULL,
  "product_id" text NULL,
  "user_id" uuid NULL,
  "product_instance_id" uuid NULL,
  PRIMARY KEY ("id")
);
-- create index "idx_user_instance" to table: "user_product_instances"
CREATE UNIQUE INDEX "idx_user_instance" ON "user_product_instances" ("user_id", "product_instance_id");
-- create "user_app_permissions" table
CREATE TABLE "user_app_permissions" (
  "id" bigserial NOT NULL,
  "user_product_instance_id" bigint NULL,
  "user_id" uuid NULL,
  "app" text NULL,
  "permission_level" text NULL,
  PRIMARY KEY ("id"),
  CONSTRAINT "fk_user_app_permissions_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT "fk_user_app_permissions_user_product_instance" FOREIGN KEY ("user_product_instance_id") REFERENCES "user_product_instances" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION
);
-- create index "idx_user_app_instance" to table: "user_app_permissions"
CREATE UNIQUE INDEX "idx_user_app_instance" ON "user_app_permissions" ("user_product_instance_id", "user_id", "app");
