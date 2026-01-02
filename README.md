# Eve Settings Manager (esm)

A CLI tool to manage Eve Online character settings. Supports listing, copying, backing up, and restoring character-specific settings across different accounts and installations.

## Features

- **List** all detected character settings with character names (via ESI API)
- **Copy** settings between characters (even across different accounts)
- **Backup** character settings to ZIP archives with metadata
- **Restore** settings from backups
- Cross-platform support (Windows & Linux/Proton)
- Works with both Steam and non-Steam installations

## Installation

### From Source

```bash
git clone https://github.com/jpbriend/eve-settings-manager.git
cd eve-settings-manager
make build
```

The binary will be in `build/esm`.

### Pre-built Binaries

Download from the [Releases](https://github.com/jpbriend/eve-settings-manager/releases) page.

## Usage

### List Characters

```bash
# List all detected character settings
esm list

# Verbose output with full paths
esm list -v
```

Output:
```
CHARACTER ID    NAME              MODIFIED
2119711681      Caldari Citizen   2024-01-15 14:30:00
95465499        Gallente Pilot    2024-01-14 09:15:00

Found 2 character(s)
```

### Backup Settings

```bash
# Backup a specific character
esm backup 2119711681

# Backup all characters
esm backup --all

# Specify output file
esm backup --all -o my-backup.zip
```

### Copy Settings

```bash
# Copy settings from one character to another
esm copy --from 2119711681 --to 95465499

# Skip confirmation prompt
esm copy --from 2119711681 --to 95465499 --force
```

This will:
1. Create a backup of the target character's settings
2. Copy the source character's settings to the target

### Restore Settings

```bash
# Restore all characters from a backup
esm restore eve-backup-20240115.zip

# Restore a specific character
esm restore eve-backup-20240115.zip -c 2119711681

# Skip confirmation
esm restore eve-backup-20240115.zip -f
```

## Settings Locations

### Windows
```
%LOCALAPPDATA%\CCP\EVE\<installation>\settings_Default\
```

### Linux (Steam/Proton)
```
~/.steam/steam/steamapps/compatdata/8500/pfx/drive_c/users/steamuser/AppData/Local/CCP/EVE/
```

## How It Works

Eve Online stores character-specific settings in `core_char_<characterID>.dat` files. These contain UI layouts, overview settings, window positions, and other character-specific preferences.

This tool:
1. Detects Eve settings directories on your system
2. Identifies character settings files
3. Resolves character IDs to names using the public ESI API
4. Allows you to manage these settings files easily

## Building

```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Install to GOPATH/bin
make install
```

## License

MIT
