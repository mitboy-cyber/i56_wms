CREATE TABLE routes (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    name VARCHAR(200),
    warehouse_id INT,
    carrier_id INT,
    transport_mode VARCHAR(30),
    origin VARCHAR(100),
    destination VARCHAR(100),
    transit_days INT,
    base_weight_price REAL DEFAULT 0,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
INSERT INTO routes (tenant_id,name,warehouse_id,carrier_id,transport_mode,origin,destination,transit_days,base_weight_price)
VALUES
    (1,'厦门→台湾(空运)',1,1,'空运','厦门','台北',3,25.00),
    (1,'深圳→台湾(海快)',2,2,'海快','深圳','基隆',5,8.00),
    (1,'厦门→台湾(海运)',1,1,'海运','厦门','高雄',8,5.00);
