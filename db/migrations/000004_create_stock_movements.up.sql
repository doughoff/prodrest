BEGIN;

DROP TYPE IF EXISTS MOVEMENT_TYPE;
CREATE TYPE MOVEMENT_TYPE AS ENUM (
    'PURCHASE',
    'ADJUST',
    'SALE',
    'PRODUCTION_OUT',
    'PRODUCTION_IN'
    );

CREATE TABLE IF NOT EXISTS "stock_movements"
(
    "id"         UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "status"     STATUS           NOT NULL DEFAULT 'ACTIVE',
    "type"       MOVEMENT_TYPE    NOT NULL,
    "date"       DATE             NOT NULL,
    "entity_id"  UUID,
    "created_by_user_id" UUID NOT NULL,
    "cancelled_by_user_id" UUID,
    "created_at" TIMESTAMP        NOT NULL DEFAULT NOW(),
    "updated_at" TIMESTAMP        NOT NULL DEFAULT NOW(),

    CONSTRAINT "fk_created_by_user"
        FOREIGN KEY("created_by_user_id")
            REFERENCES "users"("id"),

    CONSTRAINT "fk_cancelled_by_user"
        FOREIGN KEY("cancelled_by_user_id")
            REFERENCES "users"("id"),

    CONSTRAINT "fk_entity"
        FOREIGN KEY("entity_id")
            REFERENCES "entities"("id")
);

CREATE INDEX "stock_movement_status" ON "stock_movements" ("status");
CREATE INDEX "stock_movement_type" ON "stock_movements" ("type");
CREATE INDEX "stock_movement_entity_id" ON "stock_movements" ("entity_id");
CREATE INDEX "stock_movement_date" ON "stock_movements" ("date");
CREATE INDEX "stock_movement_date_month" ON "stock_movements" (date_trunc('month', "date"));

CREATE TABLE IF NOT EXISTS "stock_movement_items"
(
    "id"                UUID NOT NULL DEFAULT uuid_generate_v4(),
    "stock_movement_id" UUID             NOT NULL,
    "product_id"        UUID             NOT NULL,
    "quantity"          INT              NOT NULL,
    "price"             INT              NOT NULL,
    "batch"             varchar(30),
    "created_at"        TIMESTAMP        NOT NULL DEFAULT NOW(),
    "updated_at"        TIMESTAMP        NOT NULL DEFAULT NOW(),

    PRIMARY KEY("id", "stock_movement_id"),

    CONSTRAINT fk_stock_movement
        FOREIGN KEY(stock_movement_id)
            REFERENCES "stock_movements"(id),

    CONSTRAINT fk_product
        FOREIGN KEY(product_id)
            REFERENCES "products"(id)
);

CREATE INDEX "stock_movement_item_batch" ON "stock_movement_items" ("batch");
CREATE INDEX "stock_movement_item_product" ON "stock_movement_items" ("product_id");
CREATE INDEX "stock_movement_item_stock_movement" ON "stock_movement_items" ("stock_movement_id");

COMMIT;