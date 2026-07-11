CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    name VARCHAR(200),
    code VARCHAR(50),
    contact VARCHAR(50),
    phone VARCHAR(30),
    email VARCHAR(100),
    status VARCHAR(20) DEFAULT 'active',
    balance REAL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT NOW()
);
INSERT INTO clients (tenant_id,name,code,contact,phone,email)
VALUES
    (1,'EZ集运通','EZJYT','王经理','13900001111','contact@ezjyt.com'),
    (1,'琦立工作室','QLGZS','张总','13800003333','ql@studio.com');
