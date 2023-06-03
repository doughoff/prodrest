BEGIN;

-- Drop the index first before the table
DROP INDEX IF EXISTS "users_status";

-- Drop the table
DROP TABLE IF EXISTS "users";

-- Drop the type
DROP TYPE IF EXISTS STATUS;

COMMIT;
