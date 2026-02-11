-- Rollback retry tracking columns from source_files table

ALTER TABLE lei_raw.source_files 
    DROP COLUMN IF EXISTS retry_count,
    DROP COLUMN IF EXISTS max_retries,
    DROP COLUMN IF EXISTS failure_category;
