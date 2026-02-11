-- Add retry tracking columns to source_files table
-- This enables smart retry logic for failed file processing

ALTER TABLE lei_raw.source_files 
    ADD COLUMN retry_count INTEGER NOT NULL DEFAULT 0,
    ADD COLUMN max_retries INTEGER NOT NULL DEFAULT 3,
    ADD COLUMN failure_category VARCHAR(50);

-- Add comment to explain failure categories
COMMENT ON COLUMN lei_raw.source_files.failure_category IS 'Category of failure: SCHEMA_ERROR (retryable), NETWORK_ERROR (retryable), FILE_CORRUPTION (not retryable), UNKNOWN (retryable with caution)';
COMMENT ON COLUMN lei_raw.source_files.retry_count IS 'Number of times this file processing has been retried';
COMMENT ON COLUMN lei_raw.source_files.max_retries IS 'Maximum number of retry attempts allowed before permanent failure';
