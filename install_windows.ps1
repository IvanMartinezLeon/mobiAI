# Backup local customization files before deleting the local .pi directory
$tempTheme = Join-Path $env:TEMP "mobi-theme-bkp.json"
$tempHeader = Join-Path $env:TEMP "mobi-header-bkp.ts"
$hasBackup = $false

if ((Test-Path ".pi\themes\mobi-theme.json") -and (Test-Path ".pi\extensions\mobi-header.ts")) {
    Copy-Item -Path ".pi\themes\mobi-theme.json" -Destination $tempTheme -Force
    Copy-Item -Path ".pi\extensions\mobi-header.ts" -Destination $tempHeader -Force
    $hasBackup = $true
}

# Delete local .pi folder to prevent theme collisions
if (Test-Path ".pi") {
    Write-Host "Deleting local .pi directory to prevent theme collisions..." -ForegroundColor Yellow
    Remove-Item -Path ".pi" -Recurse -Force
}

Write-Host "=== Installing @earendil-works/pi-coding-agent for Windows ===" -ForegroundColor Cyan

# Check if Node.js is installed
$nodeCheck = Get-Command node -ErrorAction SilentlyContinue
if (-not $nodeCheck) {
    Write-Host "Error: Node.js is not installed." -ForegroundColor Red
    Write-Host "Please install Node.js before running this script." -ForegroundColor Yellow
    Write-Host "You can download it from: https://nodejs.org/" -ForegroundColor Yellow
    Write-Host "Or install it using winget: winget install OpenJS.NodeJS" -ForegroundColor Yellow
    Exit 1
}

# Check if npm is installed
$npmCheck = Get-Command npm -ErrorAction SilentlyContinue
if (-not $npmCheck) {
    Write-Host "Error: npm is not installed but Node.js was found." -ForegroundColor Red
    Write-Host "Please ensure npm is in your PATH." -ForegroundColor Yellow
    Exit 1
}

Write-Host "Node.js version: $(node -v)" -ForegroundColor Gray
Write-Host "npm version: $(npm -v)" -ForegroundColor Gray
Write-Host "Running: npm install -g --ignore-scripts @earendil-works/pi-coding-agent" -ForegroundColor Cyan

# Run installation
npm install -g --ignore-scripts @earendil-works/pi-coding-agent

if ($LASTEXITCODE -eq 0) {
    Write-Host "=== Agent installed successfully! ===" -ForegroundColor Green
    Write-Host "Installing additional Pi packages..." -ForegroundColor Cyan
    pi install npm:pi-subagents
    pi install npm:pi-mcp-adapter
    pi install npm:context-mode
    Write-Host "Installing MOBI AI theme and custom header globally..." -ForegroundColor Cyan
    
    $homeDir = [System.Environment]::GetFolderPath('UserProfile')
    $piDir = Join-Path $homeDir ".pi\agent"
    $themesDir = Join-Path $piDir "themes"
    $extensionsDir = Join-Path $piDir "extensions"
    
    New-Item -ItemType Directory -Force -Path $themesDir | Out-Null
    New-Item -ItemType Directory -Force -Path $extensionsDir | Out-Null
    
    if ($hasBackup) {
        Copy-Item -Path $tempTheme -Destination (Join-Path $themesDir "mobi-theme.json") -Force
        Copy-Item -Path $tempHeader -Destination (Join-Path $extensionsDir "mobi-header.ts") -Force
        Remove-Item -Path $tempTheme -Force
        Remove-Item -Path $tempHeader -Force
        Write-Host "Theme and extension files copied to global configuration."
        
        # Update settings.json
        $settingsPath = Join-Path $piDir "settings.json"
        $settings = @{}
        if (Test-Path $settingsPath) {
            try {
                $json = Get-Content -Raw -Path $settingsPath
                $settings = ConvertFrom-Json $json
            } catch {}
        }
        if ($null -eq $settings) { $settings = @{} }
        
        # Convert to hashtable to modify easily, or use PSCustomObject
        $settingsJson = @{}
        if ($settings -is [PSCustomObject]) {
            foreach ($prop in $settings.PSObject.Properties) {
                $settingsJson[$prop.Name] = $prop.Value
            }
        } elseif ($settings -is [System.Collections.IDictionary]) {
            $settingsJson = $settings
        }
        
        $settingsJson["theme"] = "mobi-theme"
        $settingsJson | ConvertTo-Json | Out-File -FilePath $settingsPath -Encoding utf8 -Force
        Write-Host "Global settings.json updated to use 'mobi-theme'."
    } else {
        Write-Host "Warning: Custom theme files were not found for backup." -ForegroundColor Yellow
        Write-Host "If you already have the theme installed globally, you can ignore this." -ForegroundColor Yellow
    }
    
    Write-Host "=== Installation completed successfully! ===" -ForegroundColor Green
} else {
    Write-Host "--------------------------------------------------------" -ForegroundColor Red
    Write-Host "Installation failed. If you encountered a permission error, try running PowerShell as Administrator." -ForegroundColor Yellow
    Write-Host "--------------------------------------------------------" -ForegroundColor Red
    # Clean up backups if failed
    if ($hasBackup) {
        Remove-Item -Path $tempTheme -Force
        Remove-Item -Path $tempHeader -Force
    }
    Exit $LASTEXITCODE
}
