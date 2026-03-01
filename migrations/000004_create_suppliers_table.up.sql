CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    tax_id VARCHAR(50) UNIQUE,
    email VARCHAR(255),
    phone VARCHAR(20),
    address JSONB,
    contact_person VARCHAR(100),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index para o slug pois e comum filtrarmos por slug
CREATE INDEX IF NOT EXISTS idx_suppliers_slug ON suppliers(slug);

-- Index GIN para buscas eficientes dentro do JSONB de endereço
CREATE INDEX IF NOT EXISTS idx_suppliers_address ON suppliers USING GIN (address);

-- Trigger para atualizar o campo updated_at automaticamente
CREATE TRIGGER update_suppliers_updated_at
    BEFORE UPDATE ON suppliers
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
