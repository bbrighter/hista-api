INSERT INTO "users" ("name", "password_hash_hex", "id")
VALUES ('Test', '8f84b824950f0379eef4477dcd8b7a2dbb3aceb4e1301df7d125bc552b79f851',2)
ON CONFLICT ("id")
DO UPDATE SET
    "name" = 'Test',
    "password_hash_hex" = '8f84b824950f0379eef4477dcd8b7a2dbb3aceb4e1301df7d125bc552b79f851';

SELECT setval('users_id_seq', 2);


INSERT INTO "tokens" ("bearer", "expires", "user_id", "id")
VALUES ('c79bcc0c-9467-4ef7-8b3e-c2a723c6893a', '2030-05-22 22:00:00.000', 2, 2)
ON CONFLICT ("id")
DO UPDATE SET 
    "bearer" = 'c79bcc0c-9467-4ef7-8b3e-c2a723c6893a',
    "expires" = '2030-05-22 22:00:00.000',
    "user_id" = 2;

SELECT setval('tokens_id_seq', 2);