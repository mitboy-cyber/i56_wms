-- I56 Framework v1.0 LTS — Database Migration
-- PostgreSQL 16+

-- ==========================================
-- 1. Tenants (Multi-tenant core)
-- ==========================================
CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(128) NOT NULL,
    code VARCHAR(32) NOT NULL UNIQUE,
    db_schema VARCHAR(64),          -- schema per tenant
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==========================================
-- 2. Warehouses
-- ==========================================
CREATE TABLE IF NOT EXISTS warehouses (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    code VARCHAR(16) NOT NULL,
    name VARCHAR(128) NOT NULL,
    address TEXT,
    contact_name VARCHAR(64),
    contact_phone VARCHAR(32),
    storage_daily_fee DECIMAL(10,2) DEFAULT 1.00,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- ==========================================
-- 3. Clients
-- ==========================================
CREATE TABLE IF NOT EXISTS clients (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    code VARCHAR(32) NOT NULL,
    name VARCHAR(128) NOT NULL,
    type VARCHAR(16) NOT NULL DEFAULT 'platform',
    contact_name VARCHAR(64),
    contact_phone VARCHAR(32),
    contact_email VARCHAR(128),
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    balance DECIMAL(14,2) NOT NULL DEFAULT 0,
    free_storage_days INT NOT NULL DEFAULT 30,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- ==========================================
-- 4. Client Members (子账户/收件人)
-- ==========================================
CREATE TABLE IF NOT EXISTS client_members (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id),
    name VARCHAR(64) NOT NULL,
    member_code VARCHAR(32),
    phone VARCHAR(32),
    id_type VARCHAR(32),
    id_number VARCHAR(64),
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==========================================
-- 5. Client Addresses
-- ==========================================
CREATE TABLE IF NOT EXISTS client_addresses (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id),
    member_id BIGINT REFERENCES client_members(id),
    country VARCHAR(32) NOT NULL DEFAULT 'TW',
    city VARCHAR(64),
    address TEXT NOT NULL,
    phone VARCHAR(32),
    is_default BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ==========================================
-- 6. Parcels (核心实体)
-- ==========================================
CREATE TABLE IF NOT EXISTS parcels (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    tracking_number VARCHAR(64) NOT NULL,
    courier_code VARCHAR(16),
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    client_id BIGINT REFERENCES clients(id),
    member_id BIGINT REFERENCES client_members(id),
    product_name VARCHAR(255),
    parcel_name VARCHAR(255),
    cargo_type VARCHAR(32) DEFAULT 'general',
    sea_freight_category VARCHAR(32),
    status VARCHAR(32) NOT NULL DEFAULT 'pre_declared',
    weight DECIMAL(10,3),
    length DECIMAL(10,2),
    width DECIMAL(10,2),
    height DECIMAL(10,2),
    location_barcode VARCHAR(64),
    photo_url TEXT,
    remark TEXT,
    receiver_name VARCHAR(64),
    receiver_phone VARCHAR(32),
    abnormal_reason VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, tracking_number)
);
CREATE INDEX idx_parcels_status ON parcels(tenant_id, status);
CREATE INDEX idx_parcels_client ON parcels(tenant_id, client_id);
CREATE INDEX idx_parcels_warehouse ON parcels(tenant_id, warehouse_id);

-- ==========================================
-- 7. Orders
-- ==========================================
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    order_no VARCHAR(32) NOT NULL UNIQUE,
    client_id BIGINT NOT NULL REFERENCES clients(id),
    member_id BIGINT NOT NULL REFERENCES client_members(id),
    route_id BIGINT,
    transport_type VARCHAR(16),
    cargo_type VARCHAR(32),
    tax_type VARCHAR(16),
    status VARCHAR(32) NOT NULL DEFAULT 'draft',
    total_weight DECIMAL(10,3),
    total_volume DECIMAL(10,3),
    weight_price DECIMAL(10,2),
    volume_price DECIMAL(10,2),
    service_fee DECIMAL(10,2) DEFAULT 0,
    delivery_fee DECIMAL(10,2) DEFAULT 0,
    total_amount DECIMAL(10,2),
    cost_amount DECIMAL(10,2),
    profit_amount DECIMAL(10,2),
    receiver_name VARCHAR(64),
    receiver_phone VARCHAR(32),
    receiver_address TEXT,
    carrier_id BIGINT,
    customs_broker_id BIGINT,
    customs_number VARCHAR(64),
    container_no VARCHAR(32),
    seal_no VARCHAR(32),
    tracking_no VARCHAR(64),
    paid_at TIMESTAMPTZ,
    shipped_at TIMESTAMPTZ,
    delivered_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Order-Parcel join table
CREATE TABLE IF NOT EXISTS order_parcels (
    order_id BIGINT NOT NULL REFERENCES orders(id),
    parcel_id BIGINT NOT NULL REFERENCES parcels(id),
    PRIMARY KEY(order_id, parcel_id)
);

-- ==========================================
-- 8. Routes (Pricing Templates)
-- ==========================================
CREATE TABLE IF NOT EXISTS routes (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    code VARCHAR(32) NOT NULL,
    name VARCHAR(128) NOT NULL,
    transport_type VARCHAR(16) NOT NULL,
    origin_warehouse_id BIGINT REFERENCES warehouses(id),
    area_group_id BIGINT,
    cargo_type VARCHAR(32) NOT NULL,
    tax_type VARCHAR(16) NOT NULL,
    weight_unit_price DECIMAL(10,2),
    volume_unit_price DECIMAL(10,2),
    min_charge DECIMAL(10,2) NOT NULL DEFAULT 50.00,
    first_weight DECIMAL(10,3) DEFAULT 1.0,
    first_weight_price DECIMAL(10,2),
    additional_weight_price DECIMAL(10,2),
    first_volume DECIMAL(10,3) DEFAULT 1.0,
    first_volume_price DECIMAL(10,2),
    additional_volume_price DECIMAL(10,2),
    estimated_days INT,
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- ==========================================
-- 9. Client Pricing (per-client price override)
-- ==========================================
CREATE TABLE IF NOT EXISTS client_pricing (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id),
    route_id BIGINT NOT NULL REFERENCES routes(id),
    transport_type VARCHAR(16),
    cargo_type VARCHAR(32),
    tax_type VARCHAR(16),
    weight_unit_price DECIMAL(10,2),
    volume_unit_price DECIMAL(10,2),
    min_charge DECIMAL(10,2),
    first_weight DECIMAL(10,3),
    first_weight_price DECIMAL(10,2),
    additional_weight_price DECIMAL(10,2),
    first_volume DECIMAL(10,3),
    first_volume_price DECIMAL(10,2),
    additional_volume_price DECIMAL(10,2),
    status VARCHAR(16) NOT NULL DEFAULT 'active',
    UNIQUE(client_id, route_id)
);

-- ==========================================
-- 10. Ledger (Financial Journal)
-- ==========================================
CREATE TABLE IF NOT EXISTS ledgers (
    id BIGSERIAL PRIMARY KEY,
    client_id BIGINT NOT NULL REFERENCES clients(id),
    type VARCHAR(16) NOT NULL,
    amount DECIMAL(14,2) NOT NULL,
    balance_after DECIMAL(14,2) NOT NULL,
    reference_no VARCHAR(64),
    tracking_no VARCHAR(64),
    remark TEXT,
    operator VARCHAR(64),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_ledgers_client ON ledgers(client_id, created_at DESC);

-- ==========================================
-- 11. Audit Log
-- ==========================================
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    operator VARCHAR(64) NOT NULL,
    module VARCHAR(32) NOT NULL,
    action VARCHAR(32) NOT NULL,
    target VARCHAR(255) NOT NULL,
    result VARCHAR(16) NOT NULL DEFAULT 'success',
    detail JSONB,
    ip_address VARCHAR(45),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_audit_tenant_time ON audit_logs(tenant_id, created_at DESC);

-- ==========================================
-- Seed Data
-- ==========================================
INSERT INTO tenants (id, name, code) VALUES (1, 'I56 Demo', 'i56') ON CONFLICT DO NOTHING;
INSERT INTO warehouses (id, tenant_id, code, name, address) VALUES (1, 1, 'XM', '厦门仓', '海沧区新盛路26号') ON CONFLICT DO NOTHING;
INSERT INTO clients (id, tenant_id, code, name, type) VALUES (1, 1, 'EZ001', 'EZ集运通', 'platform') ON CONFLICT DO NOTHING;
INSERT INTO client_members (id, client_id, name, member_code, phone) VALUES 
    (1, 1, '王仁照', '127518', '886912345678'),
    (2, 1, '吴欣如', '127680', '886923456789') ON CONFLICT DO NOTHING;
