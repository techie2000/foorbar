-- Rollback VARCHAR length standardization
-- Reverts to mixed sizes: 250 for registration, 200 for city/region

-- Revert registration and entity fields: 255 → 250
ALTER TABLE lei_raw.lei_records
ALTER COLUMN registration_authority TYPE VARCHAR(250),
ALTER COLUMN registration_authority_id TYPE VARCHAR(250),
ALTER COLUMN registration_number TYPE VARCHAR(250),
ALTER COLUMN entity_legal_form TYPE VARCHAR(250),
ALTER COLUMN managing_lou TYPE VARCHAR(250),
ALTER COLUMN validation_authority TYPE VARCHAR(250);

-- Revert address city/region fields: 255 → 200
ALTER TABLE lei_raw.lei_records
ALTER COLUMN legal_address_city TYPE VARCHAR(200),
ALTER COLUMN legal_address_region TYPE VARCHAR(200),
ALTER COLUMN hq_address_city TYPE VARCHAR(200),
ALTER COLUMN hq_address_region TYPE VARCHAR(200);
