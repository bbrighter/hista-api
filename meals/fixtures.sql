INSERT INTO "meals" ("date", "id") 
VALUES ('2024-05-22 20:37:37.477', 1)
ON CONFLICT ("id")
DO NOTHING;

SELECT setval('meals_id_seq', 2, true);

INSERT INTO "ingredients" ("name", "id")
VALUES ('Name',1)
ON CONFLICT ("id")
DO NOTHING;

SELECT setval('ingredients_id_seq', 2, true);

INSERT INTO "foods" ("ingredient_id","condition","meal_id", "id") 
VALUES (1,'',1,1)
ON CONFLICT ("id") 
DO NOTHING;

SELECT setval('foods_id_seq', 2, true);