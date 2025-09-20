-- reverse: create index "idx_user_app_instance" to table: "user_app_permissions"
DROP INDEX "idx_user_app_instance";
-- reverse: create "user_app_permissions" table
DROP TABLE "user_app_permissions";
-- reverse: create index "idx_user_instance" to table: "user_product_instances"
DROP INDEX "idx_user_instance";
-- reverse: create "user_product_instances" table
DROP TABLE "user_product_instances";
-- reverse: create index "idx_users_name" to table: "users"
DROP INDEX "idx_users_name";
-- reverse: create "users" table
DROP TABLE "users";
