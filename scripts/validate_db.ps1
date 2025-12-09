# Database Integrity Validation Script
# Usage: .\validate_db.ps1

Write-Host "Running database integrity validation..." -ForegroundColor Cyan

# Run the validation script from the project root
go run .\scripts\validate_db.go

# Check exit code
if ($LASTEXITCODE -eq 0) {
    Write-Host "`nValidation completed successfully!" -ForegroundColor Green
} else {
    Write-Host "`nValidation found issues. Please review the output above." -ForegroundColor Yellow
    exit 1
}
