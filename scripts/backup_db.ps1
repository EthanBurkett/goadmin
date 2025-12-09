# Database Backup Script
# Usage: .\backup_db.ps1 [-BackupDir <directory>] [-DbPath <path>]

param(
    [string]$BackupDir = "backups",
    [string]$DbPath = "data.sqlite"
)

Write-Host "Creating database backup..." -ForegroundColor Cyan

# Run the backup script
go run .\scripts\backup_db.go -dir $BackupDir -db $DbPath

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ Backup completed successfully!" -ForegroundColor Green
} else {
    Write-Host "`n❌ Backup failed!" -ForegroundColor Red
    exit 1
}
