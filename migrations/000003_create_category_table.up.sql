CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(100) NOT NULL UNIQUE,
    slug VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    parent_id UUID REFERENCES categories(id) ON DELETE SET NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index para o slug pois e comum filtrarmos produtos por slug de categoria
CREATE INDEX IF NOT EXISTS idx_categories_slug ON categories(slug);

-- Index para o parent_id para otimizar a busca de subcategorias
CREATE INDEX IF NOT EXISTS idx_categories_parent_id ON categories(parent_id);

-- Trigger para atualizar o campo updated_at automaticamente
-- A funcao update_updated_at_column() ja foi criada na migration 000002
CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON categories
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
