-- create "notes" table
CREATE TABLE "notes" (
  "id" bigserial NOT NULL,
  "date" timestamptz NULL,
  "text" text NULL,
  PRIMARY KEY ("id")
);
