-- Rollback: Remove all database documentation comments
-- Note: This removes documentation but does not affect data or schema structure

-- Remove schema comment
COMMENT ON SCHEMA lei_raw IS NULL;

-- Remove table comments - PUBLIC schema
COMMENT ON TABLE countries IS NULL;
COMMENT ON TABLE currencies IS NULL;
COMMENT ON TABLE addresses IS NULL;
COMMENT ON TABLE entities IS NULL;
COMMENT ON TABLE entity_addresses IS NULL;

-- Remove table comments - LEI_RAW schema
COMMENT ON TABLE lei_raw.source_files IS NULL;
COMMENT ON TABLE lei_raw.lei_records IS NULL;
COMMENT ON TABLE lei_raw.lei_records_audit IS NULL;
COMMENT ON TABLE lei_raw.file_processing_status IS NULL;

-- Remove all column comments (batch removal)
-- PUBLIC schema
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN 
        SELECT table_name, column_name 
        FROM information_schema.columns 
        WHERE table_schema = 'public' 
        AND table_name IN ('countries', 'currencies', 'addresses', 'entities', 'entity_addresses')
    LOOP
        EXECUTE FORMAT('COMMENT ON COLUMN public.%I.%I IS NULL', r.table_name, r.column_name);
    END LOOP;
END $$;

-- LEI_RAW schema
DO $$
DECLARE
    r RECORD;
BEGIN
    FOR r IN 
        SELECT table_name, column_name 
        FROM information_schema.columns 
        WHERE table_schema = 'lei_raw'
        AND table_name IN ('source_files', 'lei_records', 'lei_records_audit', 'file_processing_status')
    LOOP
        EXECUTE FORMAT('COMMENT ON COLUMN lei_raw.%I.%I IS NULL', r.table_name, r.column_name);
    END LOOP;
END $$;
