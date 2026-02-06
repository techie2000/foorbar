-- Create UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create countries table
CREATE TABLE IF NOT EXISTS countries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(2) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    alpha3_code VARCHAR(3),
    region VARCHAR(255),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_countries_code ON countries(code);
CREATE INDEX idx_countries_deleted_at ON countries(deleted_at);

-- Create currencies table
CREATE TABLE IF NOT EXISTS currencies (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(3) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(10),
    decimal_places INTEGER DEFAULT 2,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_currencies_code ON currencies(code);
CREATE INDEX idx_currencies_deleted_at ON currencies(deleted_at);

-- Create addresses table (ISO20022 compliant)
CREATE TABLE IF NOT EXISTS addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    -- Structured address fields (ISO20022)
    address_type VARCHAR(50),                -- Type of address (e.g., ADDR, PBOX, HOME, BIZZ)
    department VARCHAR(70),                  -- Department
    sub_department VARCHAR(70),              -- Sub-department
    street_name VARCHAR(70),                 -- Street name
    building_number VARCHAR(16),             -- Building number
    building_name VARCHAR(35),               -- Building name
    floor VARCHAR(70),                       -- Floor
    post_box VARCHAR(16),                    -- Post office box number
    room VARCHAR(70),                        -- Room
    postal_code VARCHAR(16),                 -- Postal code/ZIP code
    town_name VARCHAR(35),                   -- Town/city name
    town_location_name VARCHAR(35),          -- Town location name
    district_name VARCHAR(35),               -- District name
    country_sub_division VARCHAR(35),        -- State/province/region
    country_id UUID REFERENCES countries(id), -- Country reference
    -- Unstructured address lines (ISO20022 allows up to 7 lines)
    address_line_1 VARCHAR(70),
    address_line_2 VARCHAR(70),
    address_line_3 VARCHAR(70),
    address_line_4 VARCHAR(70),
    address_line_5 VARCHAR(70),
    address_line_6 VARCHAR(70),
    address_line_7 VARCHAR(70),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_addresses_country_id ON addresses(country_id);
CREATE INDEX idx_addresses_postal_code ON addresses(postal_code);
CREATE INDEX idx_addresses_town_name ON addresses(town_name);
CREATE INDEX idx_addresses_deleted_at ON addresses(deleted_at);

-- Create entities table
CREATE TABLE IF NOT EXISTS entities (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    registration_number VARCHAR(255) UNIQUE,
    type VARCHAR(50),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_entities_registration_number ON entities(registration_number);
CREATE INDEX idx_entities_type ON entities(type);
CREATE INDEX idx_entities_deleted_at ON entities(deleted_at);

-- Create entity_addresses junction table (many-to-many relationship)
CREATE TABLE IF NOT EXISTS entity_addresses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID NOT NULL REFERENCES entities(id) ON DELETE CASCADE,
    address_id UUID NOT NULL REFERENCES addresses(id) ON DELETE CASCADE,
    address_type VARCHAR(50),  -- e.g., 'REGISTERED', 'TRADING', 'BILLING', 'CORRESPONDENCE'
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    UNIQUE(entity_id, address_id)
);

CREATE INDEX idx_entity_addresses_entity_id ON entity_addresses(entity_id);
CREATE INDEX idx_entity_addresses_address_id ON entity_addresses(address_id);
CREATE INDEX idx_entity_addresses_deleted_at ON entity_addresses(deleted_at);

-- Create instruments table
CREATE TABLE IF NOT EXISTS instruments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    isin VARCHAR(12) UNIQUE,  -- ISIN is optional as not all instruments have one
    name VARCHAR(255) NOT NULL,
    type VARCHAR(50),
    currency_id UUID REFERENCES currencies(id),
    exchange VARCHAR(100),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_instruments_isin ON instruments(isin) WHERE isin IS NOT NULL;
CREATE INDEX idx_instruments_type ON instruments(type);
CREATE INDEX idx_instruments_currency_id ON instruments(currency_id);
CREATE INDEX idx_instruments_deleted_at ON instruments(deleted_at);

-- Create accounts table
CREATE TABLE IF NOT EXISTS accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_number VARCHAR(255) NOT NULL UNIQUE,
    entity_id UUID REFERENCES entities(id),
    currency_id UUID REFERENCES currencies(id),
    type VARCHAR(50),
    balance DECIMAL(19, 4) DEFAULT 0,
    opened_at TIMESTAMP NOT NULL DEFAULT NOW(),
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_accounts_account_number ON accounts(account_number);
CREATE INDEX idx_accounts_entity_id ON accounts(entity_id);
CREATE INDEX idx_accounts_type ON accounts(type);
CREATE INDEX idx_accounts_deleted_at ON accounts(deleted_at);

-- Create SSIs table
CREATE TABLE IF NOT EXISTS ssis (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_id UUID REFERENCES entities(id),
    currency_id UUID REFERENCES currencies(id),
    instrument_id UUID REFERENCES instruments(id),
    beneficiary_name VARCHAR(255) NOT NULL,
    beneficiary_account VARCHAR(255) NOT NULL,
    beneficiary_bank VARCHAR(255) NOT NULL,
    beneficiary_bank_bic VARCHAR(11),
    intermediary_bank VARCHAR(255),
    intermediary_bank_bic VARCHAR(11),
    settlement_type VARCHAR(50),
    valid_from TIMESTAMP NOT NULL DEFAULT NOW(),
    valid_to TIMESTAMP,
    active BOOLEAN DEFAULT true,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_ssis_entity_id ON ssis(entity_id);
CREATE INDEX idx_ssis_currency_id ON ssis(currency_id);
CREATE INDEX idx_ssis_instrument_id ON ssis(instrument_id);
CREATE INDEX idx_ssis_deleted_at ON ssis(deleted_at);

-- Create audit_logs table
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    entity_type VARCHAR(100) NOT NULL,
    entity_id UUID NOT NULL,
    action VARCHAR(50) NOT NULL,
    user_id UUID,
    changed_data JSONB,
    previous_data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE INDEX idx_audit_logs_entity_type ON audit_logs(entity_type);
CREATE INDEX idx_audit_logs_entity_id ON audit_logs(entity_id);
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_action ON audit_logs(action);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);

-- Create trigger function for updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers for all tables
CREATE TRIGGER update_countries_updated_at BEFORE UPDATE ON countries
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_currencies_updated_at BEFORE UPDATE ON currencies
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_addresses_updated_at BEFORE UPDATE ON addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_entities_updated_at BEFORE UPDATE ON entities
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_entity_addresses_updated_at BEFORE UPDATE ON entity_addresses
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_instruments_updated_at BEFORE UPDATE ON instruments
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_accounts_updated_at BEFORE UPDATE ON accounts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_ssis_updated_at BEFORE UPDATE ON ssis
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_audit_logs_updated_at BEFORE UPDATE ON audit_logs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
