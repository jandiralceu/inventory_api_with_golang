CREATE TABLE IF NOT EXISTS products (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Stock Keeping Unit - unique product code for internal tracking
    sku VARCHAR(50) NOT NULL UNIQUE,
    
    -- URL-friendly version of the product name (e.g., 'iphone-13-black')
    slug VARCHAR(200) NOT NULL UNIQUE,
    
    -- Display name of the product
    name VARCHAR(200) NOT NULL,
    
    -- Detailed product description
    description TEXT,
    
    -- Selling price (must be non-negative)
    price DECIMAL(10,2) NOT NULL CHECK (price >= 0),
    
    -- Cost price from supplier (for profit margin calculations)
    cost_price DECIMAL(10,2) CHECK (cost_price >= 0),
    
    -- Reference to the product category (hierarchical categories table)
    category_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    
    -- Reference to the main supplier
    supplier_id UUID REFERENCES suppliers(id) ON DELETE SET NULL,
    
    -- Soft delete / visibility control
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    
    -- Stock level that triggers a reorder alert (when quantity <= reorder_level)
    reorder_level INTEGER NOT NULL DEFAULT 0,
    
    -- Quantity to order when reorder level is reached
    reorder_quantity INTEGER NOT NULL DEFAULT 0,
    
    -- Product weight in kilograms (useful for shipping calculations)
    weight_kg DECIMAL(10,3),
    
    -- Array of product image URLs or image objects
    images JSONB,
    
    -- Flexible JSON field for any additional product attributes
    metadata JSONB,
    
    -- Timestamps with timezone
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_products_sku ON products(sku);
CREATE INDEX IF NOT EXISTS idx_products_slug ON products(slug);
CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_supplier_id ON products(supplier_id);

-- Index for filtering active/inactive products
CREATE INDEX IF NOT EXISTS idx_products_is_active ON products(is_active);

-- Index for price range queries
CREATE INDEX IF NOT EXISTS idx_products_price ON products(price);

-- GIN indexes for JSONB columns (allows efficient querying inside JSON)
CREATE INDEX IF NOT EXISTS idx_products_images ON products USING GIN (images);
CREATE INDEX IF NOT EXISTS idx_products_metadata ON products USING GIN (metadata);

-- Trigger to automatically update the updated_at timestamp on row updates
CREATE TRIGGER update_products_updated_at
    BEFORE UPDATE ON products
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
