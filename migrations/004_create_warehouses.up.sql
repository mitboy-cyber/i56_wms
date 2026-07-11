CREATE TABLE warehouses (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    name VARCHAR(200),
    code VARCHAR(20),
    province VARCHAR(50),
    city VARCHAR(50),
    district VARCHAR(50),
    contact VARCHAR(50),
    phone VARCHAR(30),
    address TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
INSERT INTO warehouses (tenant_id,name,code,province,city,district,contact,phone)
VALUES
    (1,'厦门仓','WH001','福建省','厦门市','海沧区','霏霏','13900000001'),
    (1,'深圳仓','WH002','广东省','深圳市','宝安区','王经理','13800002222');
