-- Standardize VARCHAR lengths to 255 for consistency and PostgreSQL optimization
-- PostgreSQL uses 1-byte length storage for VARCHAR(255) and below
-- Analysis of 2.982M records showed max actual length is 113 chars (registration_number)
-- This provides 2.25x headroom while optimizing storage

-- Registration and entity fields: 250 → 255
ALTER TABLE lei_raw.lei_records
ALTER COLUMN registration_authority TYPE VARCHAR(255),
ALTER COLUMN registration_authority_id TYPE VARCHAR(255),
ALTER COLUMN registration_number TYPE VARCHAR(255),
ALTER COLUMN entity_legal_form TYPE VARCHAR(255),
ALTER COLUMN managing_lou TYPE VARCHAR(255),
ALTER COLUMN validation_authority TYPE VARCHAR(255);

-- Address city/region fields: 200 → 255
ALTER TABLE lei_raw.lei_records
ALTER COLUMN legal_address_city TYPE VARCHAR(255),
ALTER COLUMN legal_address_region TYPE VARCHAR(255),
ALTER COLUMN hq_address_city TYPE VARCHAR(255),
ALTER COLUMN hq_address_region TYPE VARCHAR(255);
