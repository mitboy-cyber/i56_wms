-- I56 Device Gateway: Rollback
-- Migration 009: Drop device tables

DROP TABLE IF EXISTS device_registry CASCADE;
DROP TABLE IF EXISTS inbound_tasks CASCADE;
DROP TABLE IF EXISTS weight_records CASCADE;
