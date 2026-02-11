-- Rollback: Revert VARCHAR(255) columns back to VARCHAR(50)
-- WARNING: This will truncate data if values exceed 50 characters

ALTER TABLE lei_raw.lei_records
ALTER COLUMN entity_category TYPE VARCHAR(50),
ALTER COLUMN entity_status TYPE VARCHAR(50),
ALTER COLUMN entity_sub_category TYPE VARCHAR(50),
ALTER COLUMN legal_address_postal_code TYPE VARCHAR(50),
ALTER COLUMN hq_address_postal_code TYPE VARCHAR(50);
