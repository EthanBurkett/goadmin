# Plugin Auto-Import Script
# Automatically discovers and imports plugins in the GoAdmin codebase

Write-Host "GoAdmin Plugin System - Registry-Based (Windows Compatible)" -ForegroundColor Cyan
Write-Host "============================================================`n" -ForegroundColor Cyan

Write-Host "Plugin Architecture:" -ForegroundColor Yellow
Write-Host "  - Plugins are compiled directly into the GoAdmin binary"
Write-Host "  - No separate .so files needed (Windows compatible)"
Write-Host "  - Plugins self-register via init() functions"
Write-Host ""

$ProjectRoot = Split-Path -Parent $PSScriptRoot
$PluginsDir = Join-Path $ProjectRoot "plugins"
$MainGoPath = Join-Path $ProjectRoot "app\main.go"

Write-Host "Scanning for plugins in $PluginsDir..." -ForegroundColor Yellow

# Find all plugin directories (recursively, containing .go files)
$AllPluginDirs = @()
Get-ChildItem -Path $PluginsDir -Directory -Recurse | ForEach-Object {
    $GoFiles = Get-ChildItem -Path $_.FullName -Filter "*.go" -File
    if ($GoFiles.Count -gt 0) {
        $AllPluginDirs += $_
    }
}

if ($AllPluginDirs.Count -eq 0) {
    Write-Host "No plugins found in $PluginsDir" -ForegroundColor Red
    exit 0
}

Write-Host "`nFound $($AllPluginDirs.Count) plugin(s):`n" -ForegroundColor Green

# Read main.go content
if (-not (Test-Path $MainGoPath)) {
    Write-Host "Error: main.go not found at $MainGoPath" -ForegroundColor Red
    exit 1
}

$MainGoContent = Get-Content $MainGoPath -Raw
$ImportsToAdd = @()
$ActivePlugins = @()
$InactivePlugins = @()

foreach ($PluginDir in $AllPluginDirs) {
    # Calculate relative path from plugins directory
    $RelativePath = $PluginDir.FullName.Substring($PluginsDir.Length + 1).Replace('\', '/')
    $ImportPath = "github.com/ethanburkett/goadmin/plugins/$RelativePath"
    
    $GoFiles = Get-ChildItem -Path $PluginDir.FullName -Filter "*.go" -File
    
    Write-Host "  📦 $RelativePath" -ForegroundColor Cyan
    Write-Host "     Path: $($PluginDir.FullName)" -ForegroundColor Gray
    Write-Host "     Files: $($GoFiles.Count) Go file(s)" -ForegroundColor Gray
    
    if ($MainGoContent -match [regex]::Escape($ImportPath)) {
        Write-Host "     Status: ✓ Already imported (ACTIVE)" -ForegroundColor Green
        $ActivePlugins += $RelativePath
    } else {
        Write-Host "     Status: + Will be imported (ACTIVATING)" -ForegroundColor Yellow
        $ImportsToAdd += $ImportPath
        $InactivePlugins += $RelativePath
    }
    Write-Host ""
}

# Auto-import plugins if needed
if ($ImportsToAdd.Count -gt 0) {
    Write-Host "================================================" -ForegroundColor Cyan
    Write-Host "Auto-importing $($ImportsToAdd.Count) plugin(s)..." -ForegroundColor Yellow
    Write-Host ""
    
    # Find the plugin imports section in main.go
    $ImportSectionPattern = '(?s)(// Import plugins to register them\s*\n)(\s*_\s+"[^"]+"\s*\n)*'
    
    if ($MainGoContent -match $ImportSectionPattern) {
        # Build new imports section
        $NewImports = "// Import plugins to register them`n"
        
        # Get existing imports
        $ExistingImports = @()
        if ($Matches[2]) {
            $ExistingImports = ($Matches[2] -split "`n" | Where-Object { $_ -match '_\s+"([^"]+)"' } | ForEach-Object {
                if ($_ -match '"([^"]+)"') { $Matches[1] }
            })
        }
        
        # Combine and sort all imports
        $AllImports = ($ExistingImports + $ImportsToAdd) | Sort-Object | Select-Object -Unique
        
        foreach ($Import in $AllImports) {
            $NewImports += "`t_ `"$Import`"`n"
        }
        
        # Replace the import section
        $MainGoContent = $MainGoContent -replace $ImportSectionPattern, $NewImports
        
        # Write back to main.go
        Set-Content -Path $MainGoPath -Value $MainGoContent -NoNewline
        
        foreach ($Import in $ImportsToAdd) {
            Write-Host "  ✓ Added: $Import" -ForegroundColor Green
        }
        
        Write-Host "`n✓ Updated app/main.go with new plugin imports" -ForegroundColor Green
        Write-Host "`n⚠ Next steps:" -ForegroundColor Yellow
        Write-Host "  1. Rebuild: go build -o goadmin.exe ./app" -ForegroundColor White
        Write-Host "  2. Restart GoAdmin to activate plugins" -ForegroundColor White
    } else {
        Write-Host "Error: Could not find plugin import section in main.go" -ForegroundColor Red
        Write-Host "Expected pattern: // Import plugins to register them" -ForegroundColor Red
        exit 1
    }
} else {
    Write-Host "================================================" -ForegroundColor Cyan
    Write-Host "✓ All plugins are already imported!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Active plugins:" -ForegroundColor Yellow
    foreach ($Plugin in $ActivePlugins) {
        Write-Host "  - $Plugin" -ForegroundColor Green
    }
}

Write-Host "================================================" -ForegroundColor Cyan

exit 0
