#!/usr/bin/env pwsh
# Monitor LEI Full Sync Progress
# Usage: .\monitor-lei-sync.ps1

Write-Host "`n=== LEI Full Sync Progress Monitor ===" -ForegroundColor Cyan
Write-Host "Expected: 3,209,464 records (~9-10 hours processing)" -ForegroundColor Yellow
Write-Host ""

# Check download status
Write-Host "Download Status:" -ForegroundColor Green
docker logs axiom-dev-backend 2>&1 | Select-String -Pattern "downloaded successfully|Starting file processing" | Select-Object -Last 1

# Check current record count
Write-Host "`nDatabase Status:" -ForegroundColor Green
docker exec axiom-dev-postgres psql -U axiom -d axiom_dev -c "
SELECT 
    COUNT(*) as current_records,
    COALESCE((SELECT total_records FROM lei_raw.source_files WHERE file_type='FULL' ORDER BY downloaded_at DESC LIMIT 1), 0) as expected_records,
    ROUND((COUNT(*)::numeric / NULLIF((SELECT total_records FROM lei_raw.source_files WHERE file_type='FULL' ORDER BY downloaded_at DESC LIMIT 1), 0)) * 100, 2) as percent_complete,
    TO_CHAR(NOW(), 'HH24:MI:SS') as current_time
FROM lei_raw.lei_records;
" 2>&1

# Check latest processing progress from logs
Write-Host "`nLatest Processing Log:" -ForegroundColor Green
docker logs axiom-dev-backend 2>&1 | Select-String -Pattern "Processing progress" | Select-Object -Last 1

# Check source files status
Write-Host "`nSource Files:" -ForegroundColor Green
docker exec axiom-dev-postgres psql -U axiom -d axiom_dev -c "
SELECT 
    file_type,
    total_records,
    processed_records,
    failed_records,
    processing_status,
    TO_CHAR(downloaded_at, 'YYYY-MM-DD HH24:MI:SS') as downloaded
FROM lei_raw.source_files 
ORDER BY downloaded_at DESC 
LIMIT 3;
" 2>&1

Write-Host "`n=== End of Report ===" -ForegroundColor Cyan
Write-Host "Run this script again to check progress" -ForegroundColor Yellow
Write-Host ""
