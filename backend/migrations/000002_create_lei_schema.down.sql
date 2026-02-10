-- Drop LEI schema tables and indexes

DROP TRIGGER IF EXISTS update_file_processing_status_updated_at ON lei_raw.file_processing_status;
DROP TRIGGER IF EXISTS update_source_files_updated_at ON lei_raw.source_files;
DROP TRIGGER IF EXISTS update_lei_records_updated_at ON lei_raw.lei_records;

DROP TABLE IF EXISTS lei_raw.file_processing_status CASCADE;
DROP TABLE IF EXISTS lei_raw.lei_records_audit CASCADE;
DROP TABLE IF EXISTS lei_raw.lei_records CASCADE;
DROP TABLE IF EXISTS lei_raw.source_files CASCADE;

-- Drop the schema
DROP SCHEMA IF EXISTS lei_raw CASCADE;
