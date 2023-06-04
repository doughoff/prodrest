BEGIN;
DROP TYPE IF EXISTS UNIT;
CREATE TYPE UNIT AS ENUM ('KG', 'L', 'UNITS', 'OTHER');

CREATE TABLE "products"
(
    "id"                UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "status"            STATUS           NOT NULL DEFAULT 'ACTIVE',
    "name"              TEXT             NOT NULL,
    "barcode"           TEXT             NOT NULL UNIQUE,
    "unit"              UNIT             NOT NULL DEFAULT 'UNITS',
    "batch_control"     BOOLEAN          NOT NULL DEFAULT FALSE,
    "conversion_factor" NUMERIC          NOT NULL DEFAULT 1,
    "created_at"        TIMESTAMP        NOT NULL DEFAULT NOW(),
    "updated_at"        TIMESTAMP        NOT NULL DEFAULT NOW()
);

CREATE INDEX "products_status" ON "products" ("status");
COMMIT;