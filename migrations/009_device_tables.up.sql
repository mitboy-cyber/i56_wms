-- I56 Device Gateway: Database Tables
-- Migration 009: weight_records, inbound_tasks, device_registry

-- ── Weight Records (称重记录) ──
CREATE TABLE IF NOT EXISTS weight_records (
    id BIGSERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    waybill_no VARCHAR(50) NOT NULL,
    parcel_id BIGINT,
    gross_weight DECIMAL(10,3),      -- 毛重 (kg)
    tare_weight DECIMAL(10,3),       -- 皮重 (kg)
    net_weight DECIMAL(10,3),        -- 净重 (kg)
    declared_weight DECIMAL(10,3),   -- 申报重量 (kg)
    weight_diff DECIMAL(10,3),       -- 差异 (kg)
    scale_id VARCHAR(20),            -- 地磅编号
    weigh_time TIMESTAMPTZ DEFAULT NOW(),
    status SMALLINT DEFAULT 0,       -- 0:待确认 1:已确认 2:异常
    operator_id INT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_weight_records_waybill ON weight_records(waybill_no);
CREATE INDEX IF NOT EXISTS idx_weight_records_tenant ON weight_records(tenant_id);
CREATE INDEX IF NOT EXISTS idx_weight_records_scale ON weight_records(scale_id);
CREATE INDEX IF NOT EXISTS idx_weight_records_time ON weight_records(weigh_time);

-- ── Inbound Tasks (入库任务) ──
CREATE TABLE IF NOT EXISTS inbound_tasks (
    id BIGSERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    waybill_no VARCHAR(50) NOT NULL,
    tracking_number VARCHAR(100),
    sku_code VARCHAR(50),
    product_name VARCHAR(200),
    planned_qty INT DEFAULT 1,
    actual_qty INT DEFAULT 0,
    weight_id BIGINT REFERENCES weight_records(id),
    location_code VARCHAR(20),
    conveyor_id VARCHAR(20),
    status SMALLINT DEFAULT 0,       -- 0:待执行 1:执行中 2:已完成 3:异常
    err_msg TEXT,
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_inbound_tasks_tracking ON inbound_tasks(tracking_number);
CREATE INDEX IF NOT EXISTS idx_inbound_tasks_waybill ON inbound_tasks(waybill_no);
CREATE INDEX IF NOT EXISTS idx_inbound_tasks_tenant ON inbound_tasks(tenant_id);
CREATE INDEX IF NOT EXISTS idx_inbound_tasks_status ON inbound_tasks(status);

-- ── Device Registry (设备注册表) ──
CREATE TABLE IF NOT EXISTS device_registry (
    id BIGSERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    device_type VARCHAR(30),         -- scale / conveyor / scanner
    device_id VARCHAR(50) UNIQUE,
    device_name VARCHAR(100),
    protocol VARCHAR(30),            -- MODBUS_RTU / CONTINUOUS / TOLEDO
    port_config TEXT,                -- JSON: {port, baud, databits, stopbits, parity}
    warehouse_id INT,
    is_active BOOLEAN DEFAULT true,
    last_heartbeat TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_device_registry_tenant ON device_registry(tenant_id);
CREATE INDEX IF NOT EXISTS idx_device_registry_type ON device_registry(device_type);
CREATE INDEX IF NOT EXISTS idx_device_registry_device_id ON device_registry(device_id);
