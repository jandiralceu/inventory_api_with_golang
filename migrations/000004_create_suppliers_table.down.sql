DROP TRIGGER IF EXISTS update_suppliers_updated_at ON suppliers;
DROP INDEX IF EXISTS idx_suppliers_slug;
DROP INDEX IF EXISTS idx_suppliers_address;
DROP TABLE IF EXISTS suppliers;
