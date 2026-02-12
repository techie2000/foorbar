-- Rollback: Remove comment (cannot undo data changes safely)
COMMENT ON COLUMN lei_raw.source_files.failure_category IS NULL;
