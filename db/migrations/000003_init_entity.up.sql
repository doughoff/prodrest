BEGIN;

CREATE TABLE IF NOT EXISTS entities
(
    "id"         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "status"     STATUS           NOT NULL DEFAULT 'ACTIVE',
    "name"       TEXT             NOT NULL,
    "ci"         TEXT,
    "ruc"        TEXT,
    "created_at" TIMESTAMP        NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP        NOT NULL DEFAULT NOW()
);

CREATE INDEX "entities_status" ON "entities" ("status");

COMMIT;