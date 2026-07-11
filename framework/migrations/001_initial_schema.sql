-- I56 Framework 1.0 LTS - Initial Schema (PostgreSQL 16+)

CREATE TABLE IF NOT EXISTS tenants (
    id BIGSERIAL PRIMARY KEY, name VARCHAR(200) NOT NULL, code VARCHAR(50) NOT NULL UNIQUE,
    schema_name VARCHAR(50), db_name VARCHAR(100), isolation_mode VARCHAR(20) NOT NULL DEFAULT 'shared',
    is_active BOOLEAN NOT NULL DEFAULT true, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    username VARCHAR(100) NOT NULL, password_hash VARCHAR(255) NOT NULL,
    display_name VARCHAR(200), email VARCHAR(200), phone VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT true, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS roles (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    name VARCHAR(100) NOT NULL, code VARCHAR(50) NOT NULL, data_scope VARCHAR(50) NOT NULL DEFAULT 'self',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS permissions (
    id BIGSERIAL PRIMARY KEY, code VARCHAR(100) NOT NULL UNIQUE, name VARCHAR(200) NOT NULL, module VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS role_permissions (
    role_id BIGINT NOT NULL REFERENCES roles(id), permission_id BIGINT NOT NULL REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL, user_id BIGINT,
    action VARCHAR(100) NOT NULL, resource VARCHAR(100), resource_id VARCHAR(100),
    changes JSONB, ip_address VARCHAR(50), created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS clients (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    name VARCHAR(200) NOT NULL, code VARCHAR(50) NOT NULL, client_type VARCHAR(50) NOT NULL DEFAULT 'normal',
    contact_name VARCHAR(100), contact_phone VARCHAR(50), contact_email VARCHAR(200), address TEXT,
    balance DECIMAL(12,2) NOT NULL DEFAULT 0, credit_limit DECIMAL(12,2) NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS client_members (
    id BIGSERIAL PRIMARY KEY, client_id BIGINT NOT NULL REFERENCES clients(id),
    name VARCHAR(200) NOT NULL, member_code VARCHAR(50) NOT NULL, phone VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT true, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS declarants (
    id BIGSERIAL PRIMARY KEY, client_id BIGINT NOT NULL REFERENCES clients(id),
    name VARCHAR(200) NOT NULL, id_number VARCHAR(100), id_type VARCHAR(50) NOT NULL DEFAULT 'id_card',
    status VARCHAR(50) NOT NULL DEFAULT 'pending', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS member_addresses (
    id BIGSERIAL PRIMARY KEY, member_id BIGINT NOT NULL REFERENCES client_members(id),
    recipient_name VARCHAR(200) NOT NULL, phone VARCHAR(50), province VARCHAR(50), city VARCHAR(50),
    district VARCHAR(100), detail TEXT, postal_code VARCHAR(20),
    is_default BOOLEAN NOT NULL DEFAULT false, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS client_ledgers (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL, client_id BIGINT NOT NULL REFERENCES clients(id),
    amount DECIMAL(12,2) NOT NULL, balance_after DECIMAL(12,2) NOT NULL,
    type VARCHAR(50) NOT NULL, description TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS warehouses (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL REFERENCES tenants(id),
    name VARCHAR(200) NOT NULL, code VARCHAR(50) NOT NULL, address TEXT,
    contact VARCHAR(100), phone VARCHAR(50), is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS warehouse_zones (
    id BIGSERIAL PRIMARY KEY, warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    name VARCHAR(200) NOT NULL, zone_type VARCHAR(50) NOT NULL DEFAULT 'storage', code VARCHAR(50),
    is_active BOOLEAN NOT NULL DEFAULT true, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS warehouse_locations (
    id BIGSERIAL PRIMARY KEY, zone_id BIGINT NOT NULL REFERENCES warehouse_zones(id),
    code VARCHAR(50) NOT NULL, location_type VARCHAR(50) NOT NULL DEFAULT 'shelf',
    is_occupied BOOLEAN NOT NULL DEFAULT false, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS containers (
    id BIGSERIAL PRIMARY KEY, warehouse_id BIGINT NOT NULL REFERENCES warehouses(id),
    name VARCHAR(200) NOT NULL, code VARCHAR(50), container_type VARCHAR(50),
    max_weight DECIMAL(10,2), is_active BOOLEAN NOT NULL DEFAULT true, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS parcels (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL, warehouse_id BIGINT, client_id BIGINT REFERENCES clients(id),
    tracking_number VARCHAR(100), courier_code VARCHAR(20), cargo_type VARCHAR(50),
    product_name VARCHAR(200), parcel_name VARCHAR(200),
    actual_weight DECIMAL(10,3) DEFAULT 0, length DECIMAL(10,2) DEFAULT 0, width DECIMAL(10,2) DEFAULT 0, height DECIMAL(10,2) DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pre_declared', is_abnormal BOOLEAN NOT NULL DEFAULT false,
    location_code VARCHAR(50), order_id BIGINT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL, order_no VARCHAR(50) NOT NULL,
    warehouse_id BIGINT REFERENCES warehouses(id), client_id BIGINT REFERENCES clients(id),
    member_id BIGINT REFERENCES client_members(id), route_id BIGINT,
    recipient_name VARCHAR(200), tracking_numbers TEXT, parcel_count INT NOT NULL DEFAULT 0,
    status VARCHAR(50) NOT NULL DEFAULT 'pending_picking',
    total_actual_weight DECIMAL(10,3) DEFAULT 0, total_chargeable_weight DECIMAL(10,3) DEFAULT 0,
    total_price DECIMAL(12,2) DEFAULT 0, customs_number VARCHAR(100), carrier_tracking_no VARCHAR(100),
    remark TEXT, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS routes (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL, warehouse_id BIGINT REFERENCES warehouses(id),
    name VARCHAR(200) NOT NULL, transport_type VARCHAR(50) NOT NULL, area_group_id BIGINT,
    cargo_types TEXT[], min_weight DECIMAL(10,3) DEFAULT 0, max_weight DECIMAL(10,3) DEFAULT 999999,
    volume_coeff INT DEFAULT 6000, weight_rounding DECIMAL(10,3) DEFAULT 0.5,
    min_amount DECIMAL(12,2) DEFAULT 0, min_days INT DEFAULT 1, max_days INT DEFAULT 30,
    base_weight_price DECIMAL(12,2) DEFAULT 0, base_volume_price DECIMAL(12,2) DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS couriers (
    id BIGSERIAL PRIMARY KEY, name VARCHAR(200) NOT NULL, code VARCHAR(20) NOT NULL UNIQUE,
    country_region VARCHAR(100) DEFAULT '中国大陆', created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS work_orders (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL, warehouse_id BIGINT REFERENCES warehouses(id),
    title VARCHAR(200) NOT NULL, description TEXT, status VARCHAR(50) NOT NULL DEFAULT 'pending',
    priority INT NOT NULL DEFAULT 1, assigned_to BIGINT REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), completed_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS print_templates (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL, name VARCHAR(200) NOT NULL, type VARCHAR(50) NOT NULL,
    content TEXT NOT NULL, is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(), updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS webhook_subscriptions (
    id BIGSERIAL PRIMARY KEY, tenant_id BIGINT NOT NULL, client_id BIGINT REFERENCES clients(id),
    event VARCHAR(100) NOT NULL, url VARCHAR(500) NOT NULL, secret VARCHAR(200),
    is_active BOOLEAN NOT NULL DEFAULT true, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS webhook_delivery_logs (
    id BIGSERIAL PRIMARY KEY, subscription_id BIGINT REFERENCES webhook_subscriptions(id),
    event VARCHAR(100) NOT NULL, payload TEXT, status_code INT,
    error TEXT, retry_count INT NOT NULL DEFAULT 0, delivered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes
CREATE INDEX IF NOT EXISTS idx_clients_tenant ON clients(tenant_id);
CREATE INDEX IF NOT EXISTS idx_parcels_tenant ON parcels(tenant_id);
CREATE INDEX IF NOT EXISTS idx_parcels_client ON parcels(client_id);
CREATE INDEX IF NOT EXISTS idx_parcels_status ON parcels(status);
CREATE INDEX IF NOT EXISTS idx_orders_tenant ON orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_orders_status ON orders(status);
CREATE INDEX IF NOT EXISTS idx_routes_tenant ON routes(tenant_id);
CREATE INDEX IF NOT EXISTS idx_work_orders_tenant ON work_orders(tenant_id);
CREATE INDEX IF NOT EXISTS idx_webhook_subs_tenant ON webhook_subscriptions(tenant_id);
CREATE INDEX IF NOT EXISTS idx_audit_tenant ON audit_logs(tenant_id);
