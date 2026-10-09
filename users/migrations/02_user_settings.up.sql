ALTER TABLE "users" 
    ADD COLUMN "loading_mode" text NOT NULL DEFAULT 'spinner', 
    ADD COLUMN "language" text NOT NULL DEFAULT 'de-DE';
