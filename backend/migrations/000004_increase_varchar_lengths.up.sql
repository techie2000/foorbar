-- Migration: Increase VARCHAR(50) columns to VARCHAR(255) to prevent truncation errors
-- Context: Several fields were truncating data, causing import failures
-- Fields affected: entity_category, entity_status, entity_sub_category, postal codes

ALTER TABLE lei_raw.lei_records 
    ALTER COLUMN entity_category TYPE VARCHAR(255),
    ALTER COLUMN entity_status TYPE VARCHAR(255),
    ALTER COLUMN entity_sub_category TYPE VARCHAR(255),
    ALTER COLUMN legal_address_postal_code TYPE VARCHAR(255),
    ALTER COLUMN hq_address_postal_code TYPE VARCHAR(255);
