#!/bin/bash

# Exit on error
set -e

echo "=== Installing @earendil-works/pi-coding-agent for Linux ==="

# Backup local customization files before deleting the local .pi directory
TEMP_THEME="/tmp/mobi-theme-bkp.json"
TEMP_HEADER="/tmp/mobi-header-bkp.ts"
HAS_BACKUP=false

if [ -f ".pi/themes/mobi-theme.json" ] && [ -f ".pi/extensions/mobi-header.ts" ]; then
    cp ".pi/themes/mobi-theme.json" "$TEMP_THEME"
    cp ".pi/extensions/mobi-header.ts" "$TEMP_HEADER"
    HAS_BACKUP=true
fi

# Delete local .pi folder to prevent theme collisions
if [ -d ".pi" ]; then
    echo "Deleting local .pi directory to prevent theme collisions..."
    rm -rf .pi
fi

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "Error: Node.js is not installed."
    echo "Please install Node.js before running this script."
    echo "You can install it using your package manager, e.g.:"
    echo "  Debian/Ubuntu: sudo apt-get update && sudo apt-get install -y nodejs npm"
    echo "  Fedora: sudo dnf install -y nodejs npm"
    echo "  Arch Linux: sudo pacman -S nodejs npm"
    echo "Or download it from: https://nodejs.org/"
    exit 1
fi

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "Error: npm is not installed but Node.js was found."
    echo "Please install npm using your package manager."
    exit 1
fi

echo "Node.js version: $(node -v)"
echo "npm version: $(npm -v)"
echo "Running: npm install -g --ignore-scripts @earendil-works/pi-coding-agent"

# Run installation
if npm install -g --ignore-scripts @earendil-works/pi-coding-agent; then
    echo "=== Agent installed successfully! ==="
    echo "Installing MOBI AI theme and custom header globally..."
    
    # Determine the real user's home directory (even if running with sudo)
    REAL_HOME="$HOME"
    if [ -n "$SUDO_USER" ] && [ "$SUDO_USER" != "root" ]; then
        REAL_HOME=$(eval echo "~$SUDO_USER")
    fi
    
    THEMES_DIR="$REAL_HOME/.pi/agent/themes"
    EXTENSIONS_DIR="$REAL_HOME/.pi/agent/extensions"
    
    mkdir -p "$THEMES_DIR"
    mkdir -p "$EXTENSIONS_DIR"
    
    # Copy theme and extension files
    if [ "$HAS_BACKUP" = true ]; then
        cp "$TEMP_THEME" "$THEMES_DIR/mobi-theme.json"
        cp "$TEMP_HEADER" "$EXTENSIONS_DIR/mobi-header.ts"
        rm -f "$TEMP_THEME" "$TEMP_HEADER"
        echo "Theme and extension files copied to global configuration."
        
        # Update settings.json using Node.js
        node -e '
        const fs = require("fs");
        const settingsPath = process.argv[1];
        let settings = {};
        try {
          if (fs.existsSync(settingsPath)) {
            settings = JSON.parse(fs.readFileSync(settingsPath, "utf8"));
          }
        } catch (e) {}
        settings.theme = "mobi-theme";
        fs.mkdirSync(require("path").dirname(settingsPath), { recursive: true });
        fs.writeFileSync(settingsPath, JSON.stringify(settings, null, 2), "utf8");
        ' "$REAL_HOME/.pi/agent/settings.json"
        echo "Global settings.json updated to use 'mobi-theme'."
        
        # Fix permissions if run via sudo
        if [ -n "$SUDO_USER" ] && [ "$SUDO_USER" != "root" ]; then
            chown -R "$SUDO_USER" "$REAL_HOME/.pi"
        fi
    else
        echo "Warning: Custom theme files were not found for backup."
        echo "If you already have the theme installed globally, you can ignore this."
    fi
    
    echo "=== Installation completed successfully! ==="
else
    echo "--------------------------------------------------------"
    echo "Installation failed. If you encountered a permission error, try running:"
    echo "sudo npm install -g --ignore-scripts @earendil-works/pi-coding-agent"
    echo "--------------------------------------------------------"
    # Clean up backups if failed
    if [ "$HAS_BACKUP" = true ]; then
        rm -f "$TEMP_THEME" "$TEMP_HEADER"
    fi
    exit 1
fi
