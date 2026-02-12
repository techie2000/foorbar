-- Rollback: Restore VARCHAR lengths to original 100 characters
-- WARNING: This may cause data truncation if any values exceed 100 chars

ALTER TABLE lei_raw.lei_records
ALTER COLUMN registration_authority TYPE VARCHAR(100),
ALTER COLUMN registration_authority_id TYPE VARCHAR(100),
ALTER COLUMN registration_number TYPE VARCHAR(100),
ALTER COLUMN entity_legal_form TYPE VARCHAR(100),
ALTER COLUMN managing_lou TYPE VARCHAR(100),
ALTER COLUMN validation_authority TYPE VARCHAR(100),
ALTER COLUMN legal_address_city TYPE VARCHAR(100),
ALTER COLUMN legal_address_region TYPE VARCHAR(100),
ALTER COLUMN hq_address_city TYPE VARCHAR(100),
ALTER COLUMN hq_address_region TYPE VARCHAR(100);
