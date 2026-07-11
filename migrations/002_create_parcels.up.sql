CREATE TABLE parcels (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    warehouse_id INT,
    client_id INT,
    tracking_number VARCHAR(100),
    courier_code VARCHAR(20),
    cargo_type VARCHAR(50) DEFAULT 'general',
    product_name VARCHAR(500),
    length REAL DEFAULT 0,
    width REAL DEFAULT 0,
    height REAL DEFAULT 0,
    actual_weight REAL DEFAULT 0,
    declared_value REAL DEFAULT 0,
    status VARCHAR(30) DEFAULT 'pre_declared',
    remark TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);
INSERT INTO parcels (tenant_id,warehouse_id,client_id,tracking_number,cargo_type,product_name,actual_weight,status)
VALUES
    (1,1,1,'SF1234567890','普货','手机壳',0.35,'received'),
    (1,1,1,'ZTO9876543210','电子产品','蓝牙耳机',0.12,'weighed'),
    (1,1,1,'YTO1112223330','普货','T恤',0.45,'stored');
