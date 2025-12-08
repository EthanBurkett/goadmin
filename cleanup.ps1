#!/usr/bin/env pwsh
# Cleanup script to kill orphaned processes and clean temp files

Write-Host "Cleaning up GoAdmin processes and temp files..." -ForegroundColor Yellow

# Kill processes
$processes = @("air", "main")
foreach ($proc in $processes) {
    $found = Get-Process -Name $proc -ErrorAction SilentlyContinue
    if ($found) {
        Stop-Process -Name $proc -Force -ErrorAction SilentlyContinue
        Write-Host "✓ Stopped $proc processes" -ForegroundColor Green
    }
}

# Clean tmp directory
if (Test-Path "tmp") {
    Remove-Item -Path "tmp" -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "✓ Removed tmp directory" -ForegroundColor Green
}

Write-Host "`nCleanup complete!" -ForegroundColor Green
