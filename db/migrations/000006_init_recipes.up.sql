begin;


CREATE TABLE IF NOT EXISTS "recipes"
(
    "recipe_id"                 UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "recipe_group_id"           UUID             NOT NULL DEFAULT uuid_generate_v4(),
    "name"                      TEXT             NOT NULL,
    "status"                    STATUS           NOT NULL DEFAULT 'ACTIVE',
    "revision"                  INT              NOT NULL DEFAULT 1,
    "is_current"                BOOLEAN          NOT NULL DEFAULT TRUE,
    "created_by_user_id"        UUID             NOT NULL,
    "created_at"                TIMESTAMP        NOT NULL DEFAULT NOW(),

    CONSTRAINT "fk_created_by_user"
        FOREIGN KEY("created_by_user_id")
            REFERENCES "users"("id")
);

CREATE INDEX "recipe_status" ON "recipes" ("status");
CREATE INDEX "recipe_group_id" ON "recipes" ("recipe_group_id");
CREATE INDEX "recipe_is_current" ON "recipes" ("is_current");

CREATE TABLE IF NOT EXISTS "recipe_ingredients"
(
    "id"                UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    "recipe_id"         UUID             NOT NULL,
    "product_id"        UUID             NOT NULL,
    "quantity"          INTEGER          NOT NULL,
    CONSTRAINT "fk_recipe_ingredient_recipe"
        FOREIGN KEY("recipe_id")
            REFERENCES "recipes"("recipe_id"),
    CONSTRAINT "fk_recipe_ingredient_product"
        FOREIGN KEY("product_id")
            REFERENCES "products"("id")
);


commit;