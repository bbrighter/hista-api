-- create "headaches" table
CREATE TABLE "headaches" (
  "id" bigserial NOT NULL,
  "date" timestamptz NULL,
  "severity" smallint NULL,
  "types" json NULL,
  "positions" json NULL,
  "symptoms" json NULL,
  PRIMARY KEY ("id")
);
