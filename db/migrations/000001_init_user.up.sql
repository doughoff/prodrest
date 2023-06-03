BEGIN;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

DROP TYPE IF EXISTS STATUS;
CREATE  TYPE STATUS AS ENUM ('ACTIVE', 'INACTIVE');

CREATE TABLE "users"
(
    "id"         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "status"     STATUS           NOT NULL DEFAULT 'ACTIVE',
    "email"      TEXT             NOT NULL UNIQUE,
    "name"       TEXT             NOT NULL,
    "password"   TEXT             NOT NULL,
    "roles"      TEXT[]           NOT NULL DEFAULT '{}',
    "created_at" TIMESTAMP        NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP        NOT NULL DEFAULT NOW()
);

CREATE INDEX "users_status" ON "users" ("status");

COMMIT;