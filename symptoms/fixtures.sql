INSERT INTO "condition_events" ("date", "id")
VALUES ('2024-05-22 20:37:37.477', 1)
ON CONFLICT ("id")
DO UPDATE SET 
    "date"='2024-05-22 20:37:37.477';

SELECT setval('condition_events_id_seq',2,true);

INSERT INTO "symptom_categories" ("id", "name")
VALUES (1, 'Category')
ON CONFLICT ("id")
DO UPDATE SET
    "name" = 'Category';

SELECT setval('symptom_categories_id_seq',2,true);

INSERT INTO "symptoms" ("id", "name", "symptom_category_id")
VALUES (1, 'Symptom', 1)
ON CONFLICT ("id")
DO UPDATE SET 
    "name" = 'Symptom';

SELECT setval('symptoms_id_seq',2,true);

INSERT INTO "conditions" ("id", "severity", "symptom_id", "condition_event_id")
VALUES (1, 4, 1, 1)
ON CONFLICT ("id")
DO UPDATE SET 
    "severity" = 4;

SELECT setval('conditions_id_seq',2,true);
