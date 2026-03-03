CREATE TABLE IF NOT EXISTS inventory (
    -- Primary identifier
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Reference to the product and warehouse
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    warehouse_id UUID NOT NULL REFERENCES warehouses(id) ON DELETE CASCADE,
    
    -- Current physical quantity in stock
    quantity INTEGER NOT NULL DEFAULT 0 CHECK (quantity >= 0),
    
    -- Reserved quantity (e.g., items in pending orders)
    reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
    
    -- Specific location within the warehouse (aisle, rack, bin)
    location_code VARCHAR(50),
    
    -- Specific levels for this warehouse
    min_quantity INTEGER DEFAULT 0 CHECK (min_quantity >= 0),
    max_quantity INTEGER,
    
    -- Version for optimistic locking
    version INTEGER NOT NULL DEFAULT 1,
    
    -- Last time stock was physically counted
    last_counted_at TIMESTAMPTZ,
    
    -- Additional metadata (batch numbers, expiry dates, etc.)
    metadata JSONB,
    
    -- Timestamps
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    -- Constraints
    CONSTRAINT uk_inventory_product_warehouse UNIQUE(product_id, warehouse_id),
    CONSTRAINT ck_inventory_quantity_available CHECK (quantity >= reserved_quantity)
);

-- Indexes
-- product_id is already indexed by the UNIQUE constraint (first column)
CREATE INDEX IF NOT EXISTS idx_inventory_warehouse_id ON inventory(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inventory_quantity ON inventory(quantity);
CREATE INDEX IF NOT EXISTS idx_inventory_location ON inventory(location_code);
CREATE INDEX IF NOT EXISTS idx_inventory_metadata ON inventory USING GIN (metadata);

-- Trigger for updated_at
CREATE TRIGGER update_inventory_updated_at
    BEFORE UPDATE ON inventory
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
