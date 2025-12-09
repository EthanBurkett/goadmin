# Database Restore Script
# Usage: .\restore_db.ps1 -BackupFile <backup.zip> [-DbPath <path>] [-Force]

param(
    [Parameter(Mandatory=$true)]
    [string]$BackupFile,
    [string]$DbPath = "data.sqlite",
    [switch]$Force
)

Write-Host "Restoring database from backup..." -ForegroundColor Cyan

# Build command arguments
$args = @("-backup", $BackupFile, "-db", $DbPath)
if ($Force) {
    $args += "-force"
}

# Run the restore script
go run .\scripts\restore_db.go @args

if ($LASTEXITCODE -eq 0) {
    Write-Host "`n✅ Database restored successfully!" -ForegroundColor Green
} else {
    Write-Host "`n❌ Restore failed!" -ForegroundColor Red
    exit 1
}
