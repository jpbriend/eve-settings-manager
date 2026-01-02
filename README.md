# Eve Settings Manager (esm)

A simple tool to manage your Eve Online character settings. Copy your UI layout, overview settings, and window positions between characters - even across different accounts.

## Why Use This Tool?

Have you ever spent hours perfecting your Eve Online UI layout, overview settings, and window positions on one character, only to wish you could use the same setup on your alts? This tool solves that problem.

**Common use cases:**

- Copy your main character's UI settings to a new alt
- Backup your settings before a major Eve update
- Share your perfect overview setup across all your characters
- Restore settings after reinstalling Eve or switching computers

## What Are Character Settings?

Eve Online stores each character's preferences in separate files. These include:

- Window positions and sizes
- Overview settings and tabs
- UI layout and scaling
- Keyboard shortcuts
- Chat channel settings
- And more...

This tool helps you manage these files without manually digging through folders.

## Installation

### Download (Recommended)

1. Go to the [Releases](https://github.com/jpbriend/eve-settings-manager/releases) page
2. Download the file for your system:
   - **Windows**: `esm_x.x.x_windows_amd64.zip`
   - **Linux**: `esm_x.x.x_linux_amd64.tar.gz`
3. Extract the archive
4. Run `esm` from the command line (see Usage below)

### Windows Quick Start

1. Download and extract the ZIP file
2. Open Command Prompt or PowerShell
3. Navigate to the extracted folder: `cd Downloads\esm_x.x.x_windows_amd64`
4. Run commands like: `.\esm.exe list`

### Linux Quick Start

1. Download and extract the archive
2. Open a terminal
3. Make it executable: `chmod +x esm`
4. Run commands like: `./esm list`

## Usage

**Important:** Close Eve Online before using this tool. Settings changes won't take effect while the game is running.

### Step 1: List Your Characters

First, see which characters the tool can find:

```bash
esm list
```

Example output:
```
CHARACTER ID    NAME              MODIFIED
2119711681      Mudak FendLaBise  2024-01-15 14:30:00
95465499        Samia FendLaBise  2024-01-14 09:15:00

Found 2 character(s)
```

The tool automatically looks up character names using Eve's public API.

### Step 2: Backup Your Settings (Recommended)

Before making any changes, create a backup:

```bash
# Backup all characters
esm backup --all

# Or backup a specific character (by name or ID)
esm backup "Mudak FendLaBise"
esm backup 2119711681

# Save to a specific file
esm backup --all -o my-eve-backup.zip
```

This creates a ZIP file containing your settings and a metadata file with character names and timestamps.

### Step 3: Copy Settings Between Characters

Copy settings from one character to another:

```bash
# Using character names
esm copy --from "Mudak FendLaBise" --to "Samia FendLaBise"

# Using character IDs
esm copy --from 2119711681 --to 95465499
```

The tool will:
1. Show you what it's about to do
2. Ask for confirmation
3. Automatically backup the target character's settings
4. Copy the settings

Use `--force` to skip the confirmation prompt:

```bash
esm copy --from "Mudak FendLaBise" --to "Samia FendLaBise" --force
```

### Step 4: Restore Settings (If Needed)

Restore settings from a backup:

```bash
# Restore all characters from a backup
esm restore my-eve-backup.zip

# Restore only a specific character
esm restore my-eve-backup.zip -c "Mudak FendLaBise"

# Skip confirmation prompt
esm restore my-eve-backup.zip --force
```

## Command Reference

| Command | Description |
|---------|-------------|
| `esm list` | Show all detected characters |
| `esm list -v` | Show characters with full file paths |
| `esm backup <character>` | Backup one character's settings |
| `esm backup --all` | Backup all characters |
| `esm backup --all -o file.zip` | Backup to a specific file |
| `esm copy --from X --to Y` | Copy settings from X to Y |
| `esm copy --from X --to Y -f` | Copy without confirmation |
| `esm restore file.zip` | Restore all characters from backup |
| `esm restore file.zip -c X` | Restore only character X |

## Supported Platforms

| Platform | Installation Type | Status |
|----------|------------------|--------|
| Windows | Standard (non-Steam) | Supported |
| Windows | Steam | Supported |
| Linux | Steam (Proton) | Supported |
| Linux | Lutris | Supported |
| macOS | - | Not supported |

## Frequently Asked Questions

### Where does Eve store my settings?

**Windows:**
```
%LOCALAPPDATA%\CCP\EVE\<installation>\settings_Default\
```
For example: `C:\Users\YourName\AppData\Local\CCP\EVE\...\settings_Default\`

**Linux (Steam/Proton):**
```
~/.steam/steam/steamapps/compatdata/8500/pfx/drive_c/users/steamuser/AppData/Local/CCP/EVE/
```

### Can I copy settings between accounts?

Yes! This tool works across different accounts. As long as both characters have logged in at least once on your computer, you can copy settings between them.

### Will this affect my in-game skills or assets?

No. This tool only manages local UI settings files. It cannot access or modify anything stored on Eve's servers (skills, assets, wallet, etc.).

### The tool doesn't find my characters

Make sure you have:
1. Logged into Eve Online with each character at least once
2. Closed Eve Online before running the tool
3. Run the tool from the correct user account

Use `esm list -v` to see which directories the tool is searching.

### Something went wrong, how do I restore my settings?

The tool automatically creates backups before making changes. Look for `.bak` files or `.zip` backups in your Eve settings folder. You can also use `esm restore` with any backup you've created.

## Troubleshooting

### "No Eve Online settings directories found"

The tool couldn't find your Eve installation. This can happen if:
- Eve is installed in a non-standard location
- You're running the tool as a different user than the one who plays Eve
- Eve has never been run on this computer

### "Character not found"

The character name might be spelled incorrectly, or the character has never logged in on this computer. Use `esm list` to see available characters.

### Settings didn't apply in-game

Make sure Eve Online was completely closed when you ran the tool. Eve loads settings when it starts and saves them when it closes - running the tool while Eve is open won't work.

## Building from Source

For developers who want to build from source:

```bash
git clone https://github.com/jpbriend/eve-settings-manager.git
cd eve-settings-manager
make build
```

Requirements:
- Go 1.21 or later
- Make

## License

MIT License - feel free to use, modify, and distribute this tool.

## Contributing

Found a bug or have a feature request? Please open an issue on [GitHub](https://github.com/jpbriend/eve-settings-manager/issues).
