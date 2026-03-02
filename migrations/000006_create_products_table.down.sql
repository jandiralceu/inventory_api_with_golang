DROP TRIGGER IF EXISTS update_products_updated_at ON products;

DROP INDEX IF EXISTS idx_products_metadata;
DROP INDEX IF EXISTS idx_products_images;
DROP INDEX IF EXISTS idx_products_price;
DROP INDEX IF EXISTS idx_products_is_active;
DROP INDEX IF EXISTS idx_products_supplier_id;
DROP INDEX IF EXISTS idx_products_category_id;
DROP INDEX IF EXISTS idx_products_name;
DROP INDEX IF EXISTS idx_products_slug;
DROP INDEX IF EXISTS idx_products_sku;

DROP TABLE IF EXISTS products;
