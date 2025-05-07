DROP TABLE "evening_statuses";
DROP TABlE "morning_statuses";

-- modify "statuses" table
ALTER TABLE "statuses" ADD COLUMN "morning_fitness" bigint NULL, ADD COLUMN "evening_fitness" bigint NULL, ADD COLUMN "morning_sleep" bigint NULL;
