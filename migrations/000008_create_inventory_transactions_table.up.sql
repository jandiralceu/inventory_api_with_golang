CREATE TABLE IF NOT EXISTS inventory_transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    
    -- Main references
    inventory_id UUID NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    product_id UUID NOT NULL REFERENCES products(id),
    warehouse_id UUID NOT NULL REFERENCES warehouses(id),
    
    -- Who performed the operation
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    
    -- Movement details
    quantity_change INTEGER NOT NULL,
    quantity_balance INTEGER NOT NULL, -- Stock level AFTER this change
    
    -- Classification and external reference
    transaction_type VARCHAR(50) NOT NULL, -- e.g., 'IN', 'OUT', 'ADJUSTMENT'
    reference_id VARCHAR(100),             -- Optional ID of an Order, Purchase or Transfer
    
    reason TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Performance indexes for reporting
CREATE INDEX IF NOT EXISTS idx_inv_trans_inventory_id ON inventory_transactions(inventory_id);
CREATE INDEX IF NOT EXISTS idx_inv_trans_product_id ON inventory_transactions(product_id);
CREATE INDEX IF NOT EXISTS idx_inv_trans_warehouse_id ON inventory_transactions(warehouse_id);
CREATE INDEX IF NOT EXISTS idx_inv_trans_created_at ON inventory_transactions(created_at DESC);
