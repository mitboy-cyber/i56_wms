CREATE TABLE carriers (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    name VARCHAR(200),
    code VARCHAR(50),
    type VARCHAR(30),
    contact VARCHAR(50),
    phone VARCHAR(30),
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
INSERT INTO carriers (tenant_id,name,code,type,contact,phone)
VALUES
    (1,'新竹物流','HCT','宅配','陈经理','13700002222'),
    (1,'黑猫宅急便','YAMATO','宅配','李经理','13800001111');
