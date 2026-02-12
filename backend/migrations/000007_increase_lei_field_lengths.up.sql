-- Migration: Increase VARCHAR(100) columns to prevent data truncation errors
-- Context: GLEIF LEI data contains fields exceeding 100 characters
-- Error: "value too long for type character varying(100)"
-- Batch 700-800 failure indicates real-world data exceeds current limits

-- Registration and authority fields (increase to 250)
ALTER TABLE lei_raw.lei_records
ALTER COLUMN registration_authority TYPE VARCHAR(250),
ALTER COLUMN registration_authority_id TYPE VARCHAR(250),
ALTER COLUMN registration_number TYPE VARCHAR(250),
ALTER COLUMN entity_legal_form TYPE VARCHAR(250),
ALTER COLUMN managing_lou TYPE VARCHAR(250),
ALTER COLUMN validation_authority TYPE VARCHAR(250);

-- Address fields - cities and regions (increase to 200 for international names)
ALTER TABLE lei_raw.lei_records
ALTER COLUMN legal_address_city TYPE VARCHAR(200),
ALTER COLUMN legal_address_region TYPE VARCHAR(200),
ALTER COLUMN hq_address_city TYPE VARCHAR(200),
ALTER COLUMN hq_address_region TYPE VARCHAR(200);

-- Audit fields (keep at 100 - 'system' is only current value)
-- created_by and updated_by remain VARCHAR(100)
