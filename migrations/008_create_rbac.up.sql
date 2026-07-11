CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    name VARCHAR(100),
    code VARCHAR(50),
    created_at TIMESTAMPTZ DEFAULT NOW()
);
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    tenant_id INT NOT NULL REFERENCES tenants(id),
    username VARCHAR(100) UNIQUE,
    password_hash VARCHAR(200),
    display_name VARCHAR(100),
    role_id INT,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMPTZ DEFAULT NOW()
);
INSERT INTO roles (tenant_id,name,code) VALUES (1,'超级管理员','super_admin'),(1,'仓库操作员','warehouse_op');
INSERT INTO users (tenant_id,username,password_hash,display_name,role_id) VALUES (1,'admin','$2a$10$placeholder_hash','系统管理员',1),(1,'op001','$2a$10$placeholder_hash','张三',2);
