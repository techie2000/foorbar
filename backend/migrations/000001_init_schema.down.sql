-- Drop triggers
DROP TRIGGER IF EXISTS update_audit_logs_updated_at ON audit_logs;
DROP TRIGGER IF EXISTS update_ssis_updated_at ON ssis;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON accounts;
DROP TRIGGER IF EXISTS update_instruments_updated_at ON instruments;
DROP TRIGGER IF EXISTS update_entities_updated_at ON entities;
DROP TRIGGER IF EXISTS update_addresses_updated_at ON addresses;
DROP TRIGGER IF EXISTS update_currencies_updated_at ON currencies;
DROP TRIGGER IF EXISTS update_countries_updated_at ON countries;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables
DROP TABLE IF EXISTS audit_logs;
DROP TABLE IF EXISTS ssis;
DROP TABLE IF EXISTS accounts;
DROP TABLE IF EXISTS instruments;
DROP TABLE IF EXISTS entities;
DROP TABLE IF EXISTS addresses;
DROP TABLE IF EXISTS currencies;
DROP TABLE IF EXISTS countries;

-- Drop extension
DROP EXTENSION IF EXISTS "uuid-ossp";
