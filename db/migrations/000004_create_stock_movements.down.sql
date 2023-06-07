BEGIN;

-- Drop the stock_movement_items table first because it references stock_movements
DROP TABLE IF EXISTS "stock_movement_items";

-- Now that there are no more references, we can drop the stock_movements table
DROP TABLE IF EXISTS "stock_movements";

-- Drop the MOVEMENT_TYPE type
DROP TYPE IF EXISTS MOVEMENT_TYPE;

COMMIT;
