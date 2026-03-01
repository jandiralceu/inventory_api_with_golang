DROP TRIGGER IF EXISTS update_warehouses_updated_at ON warehouses;
DROP INDEX IF EXISTS idx_warehouses_slug;
DROP INDEX IF EXISTS idx_warehouses_code;
DROP INDEX IF EXISTS idx_warehouses_is_active;
DROP INDEX IF EXISTS idx_warehouses_address;
DROP TABLE IF EXISTS warehouses;
