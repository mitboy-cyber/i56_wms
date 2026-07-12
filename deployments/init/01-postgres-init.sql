-- I56 Framework v2.3 — PostgreSQL Init Script
-- Creates database schema and applies migrations

-- Enable extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- ============================================================
-- Core tables
-- ============================================================

-- Tenants
CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL UNIQUE,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Users
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    username VARCHAR(100) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    real_name VARCHAR(200),
    email VARCHAR(255),
    phone VARCHAR(50),
    role_id BIGINT DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, username)
);

-- Roles
CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, slug)
);

-- Permissions
CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL UNIQUE,
    resource VARCHAR(100),
    action VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- Role permissions (M:N)
CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    permission_id BIGINT NOT NULL REFERENCES permissions(id) ON DELETE CASCADE,
    PRIMARY KEY (role_id, permission_id)
);

-- Warehouses
CREATE TABLE IF NOT EXISTS warehouses (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    address TEXT,
    contact VARCHAR(200),
    phone VARCHAR(50),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- Clients
CREATE TABLE IF NOT EXISTS clients (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    name VARCHAR(255) NOT NULL,
    code VARCHAR(50) NOT NULL,
    client_type VARCHAR(50) DEFAULT 'platform',
    contact_name VARCHAR(200),
    contact_phone VARCHAR(50),
    contact_email VARCHAR(255),
    balance NUMERIC(12,2) DEFAULT 0,
    credit_limit NUMERIC(12,2) DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, code)
);

-- Parcels
CREATE TABLE IF NOT EXISTS parcels (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    client_id BIGINT NOT NULL REFERENCES clients(id),
    tracking_number VARCHAR(100) NOT NULL,
    product_name VARCHAR(255),
    parcel_name VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pre_declared',
    courier_code VARCHAR(50),
    cargo_type VARCHAR(50) DEFAULT 'general',
    actual_weight NUMERIC(10,3) DEFAULT 0,
    length NUMERIC(10,2) DEFAULT 0,
    width NUMERIC(10,2) DEFAULT 0,
    height NUMERIC(10,2) DEFAULT 0,
    location_code VARCHAR(100),
    container_no VARCHAR(100),
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, tracking_number)
);

-- Orders
CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    client_id BIGINT NOT NULL REFERENCES clients(id),
    order_no VARCHAR(100) NOT NULL UNIQUE,
    member_id BIGINT DEFAULT 1,
    route_id BIGINT DEFAULT 1,
    recipient_name VARCHAR(255),
    tracking_numbers TEXT,
    status VARCHAR(50) DEFAULT 'pending_picking',
    parcel_count INT DEFAULT 0,
    total_actual_weight NUMERIC(12,3) DEFAULT 0,
    total_chargeable_weight NUMERIC(12,3) DEFAULT 0,
    total_price NUMERIC(12,2) DEFAULT 0,
    carrier_tracking_no VARCHAR(200),
    customs_number VARCHAR(200),
    remark TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

-- Audit logs
CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT REFERENCES tenants(id),
    user_id BIGINT REFERENCES users(id),
    action VARCHAR(50) NOT NULL,
    resource VARCHAR(100) NOT NULL,
    resource_id VARCHAR(200),
    detail JSONB,
    ip VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

-- System config
CREATE TABLE IF NOT EXISTS system_configs (
    id BIGSERIAL PRIMARY KEY,
    tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    config_key VARCHAR(200) NOT NULL,
    config_value TEXT,
    config_type VARCHAR(50) DEFAULT 'string',
    description TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE(tenant_id, config_key)
);

-- ============================================================
-- Seed data
-- ============================================================

-- Default tenant
INSERT INTO tenants (id, name, code) VALUES (1, '默认租户', 'default')
ON CONFLICT (code) DO NOTHING;

-- Default admin user (password: admin)
INSERT INTO users (tenant_id, username, password_hash, real_name, role_id, is_active)
VALUES (1, 'admin', '$2a$10$placeholder_hash_for_admin', '系统管理员', 1, true)
ON CONFLICT (tenant_id, username) DO NOTHING;

-- Default role
INSERT INTO roles (tenant_id, name, slug, description, is_active)
VALUES (1, '超级管理员', 'super_admin', '系统最高权限', true)
ON CONFLICT (tenant_id, slug) DO NOTHING;

-- Default warehouse
INSERT INTO warehouses (tenant_id, name, code, address, contact, phone)
VALUES (1, '厦门仓', 'XM', '福建省厦门市集美区', '仓库管理员', '0592-1234567')
ON CONFLICT (tenant_id, code) DO NOTHING;

-- Default client
INSERT INTO clients (tenant_id, name, code, client_type, contact_name, balance, credit_limit)
VALUES (1, 'EZ集运通', 'EZ001', 'platform', '运营经理', 10000, 20000)
ON CONFLICT (tenant_id, code) DO NOTHING;

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_audit_logs_action ON audit_logs(action);
CREATE INDEX IF NOT EXISTS idx_audit_logs_created_at ON audit_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_audit_logs_tenant ON audit_logs(tenant_id);
CREATE INDEX IF NOT EXISTS idx_parcels_status ON parcels(status);
CREATE INDEX IF NOT EXISTS idx_parcels_tracking ON parcels(tracking_number);
CREATE INDEX IF NOT EXISTS idx_orders_order_no ON orders(order_no);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_orders_created_at ON orders(created_at DESC);
