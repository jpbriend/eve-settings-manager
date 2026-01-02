# Eve Online Settings Manager - Implementation Plan

## Overview
A Go CLI application to manage Eve Online **character-specific settings** for both Steam and non-Steam versions of the game on Windows and Linux.

## Eve Online Settings Structure

### Settings Locations (Verified)

**Windows:**
- Base path: `%LOCALAPPDATA%\CCP\EVE\`
- Full path varies by installation:
  - `%LOCALAPPDATA%\CCP\EVE\c_eve_sharedcache_tq_tranquility\settings_Default\`
  - `%LOCALAPPDATA%\CCP\EVE\d_steamlibrary_steamapps_common_eve_online_sharedcache_tq_tranquility\settings_Default\`
  - `%LOCALAPPDATA%\CCP\EVE\c_program_files_eve_sharedcache_tq_tranquility\settings_Default\`

**Linux (Proton/Steam):**
- `~/.steam/steam/steamapps/compatdata/8500/pfx/drive_c/users/steamuser/AppData/Local/CCP/EVE/c_ccp_eve_tq_tranquility/settings_Default/`

**Linux (Lutris):**
- `[lutris dir]/drive_c/users/[user]/AppData/Local/CCP/EVE/c_ccp_eve_online_tq_tranquility/settings_Default/`

### Settings Files (Character-Specific Focus)
- `core_char_<characterID>.dat` - Character-specific settings (UI layout, overview, window positions, etc.)
- `core_user_<userID>.dat` - Account-wide settings (NOT in scope)
- The character ID is the Eve Online character ID (numeric)

## Project Structure

```
eve-settings-manager/
├── cmd/
│   └── esm/
│       └── main.go           # Entry point
├── internal/
│   ├── eve/
│   │   ├── detector.go       # Detect Eve installations & settings dirs
│   │   ├── character.go      # Character settings operations
│   │   └── paths.go          # Platform-specific path resolution
│   ├── esi/
│   │   └── client.go         # ESI API client (character name lookup)
│   ├── commands/
│   │   ├── list.go           # List command
│   │   ├── copy.go           # Copy command
│   │   ├── backup.go         # Backup command
│   │   └── restore.go        # Restore command
│   └── backup/
│       └── zip.go            # ZIP archive operations
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

## CLI Commands

### `esm list`
List all detected Eve installations and character settings files.
```
esm list [--verbose]
```
Output example:
```
CHARACTER ID    NAME              MODIFIED              PATH
2119711681      Caldari Citizen   2024-01-15 14:30:00   ~/.steam/.../core_char_2119711681.dat
95465499        Gallente Pilot    2024-01-14 09:15:00   ~/.steam/.../core_char_95465499.dat
```

### `esm backup`
Create a ZIP backup of character settings.
```
esm backup <character>           # by ID or name
esm backup --all --output <path>
```
- Accepts character ID or name (resolved via ESI API)
- ZIP format with metadata (includes character names)

### `esm copy`
Copy character settings from one character to another (including cross-account).
```
esm copy --from <character> --to <character> [--force]
```
- Accepts character ID or name
- Works across different accounts
- `--force` to overwrite without confirmation
- Auto-backup of target before overwriting

### `esm restore`
Restore character settings from a backup.
```
esm restore <backup.zip> [--character <name-or-id>]
```

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `archive/zip` (stdlib) | ZIP backup creation/extraction |
| `net/http` (stdlib) | ESI API calls for character name resolution |

## ESI API Integration

**Endpoint**: `GET https://esi.evetech.net/latest/characters/{character_id}/`
- Public endpoint, no authentication required
- Returns character name, corporation, etc.
- Cached for 7 days server-side

**Use cases**:
- `esm list`: Show character names alongside IDs
- `esm copy/backup`: Accept character names as arguments (resolve to ID)
- Backup metadata: Include character names for readability

**Implementation**:
- Simple HTTP GET with JSON parsing
- Local cache (optional) to reduce API calls
- Graceful fallback to ID-only if API unavailable

## Implementation Phases

### Phase 1: Core Infrastructure
1. Initialize Go module (`github.com/<user>/eve-settings-manager`)
2. Set up Cobra CLI structure with root command
3. Implement platform detection (Windows vs Linux)
4. Implement Eve settings directory discovery:
   - Scan `%LOCALAPPDATA%\CCP\EVE\` on Windows
   - Scan `~/.steam/steam/steamapps/compatdata/8500/pfx/...` on Linux
   - Find all `settings_*` directories
   - Enumerate `core_char_*.dat` files

### Phase 2: List Command
1. Scan all detected settings directories
2. Parse `core_char_<id>.dat` filenames to extract character IDs
3. Display character ID, file path, and modification time
4. `--verbose` flag for additional details

### Phase 3: Backup Command
1. Create ZIP archive using stdlib `archive/zip`
2. Include metadata JSON (timestamp, source path, character IDs)
3. Support single character (`esm backup <id>`) or all (`--all`)
4. Default output: `eve-backup-<timestamp>.zip`

### Phase 4: Copy Command
1. Locate source `core_char_<from>.dat` file
2. Locate or create target `core_char_<to>.dat` file
3. Copy file contents (preserving target filename)
4. Confirmation prompt unless `--force` specified
5. Create automatic backup of target before overwriting

### Phase 5: Restore Command
1. Extract and validate ZIP backup
2. Read metadata to identify contained characters
3. Restore to original or specified location
4. Support restoring specific character from multi-char backup

### Phase 6: Testing & Polish
1. Unit tests for path detection logic
2. Integration tests for backup/restore cycle
3. Cross-platform build via Makefile
4. README with usage examples

## Sources

- [EVE University Wiki - Client Preferences and Settings Backup](https://wiki.eveuniversity.org/Client_Preferences_and_Settings_Backup)
- [EVE Forums - Manually copy settings between characters](https://forums.eveonline.com/t/manually-copy-settings-between-characters-and-accounts/32704)
- [EVE University Wiki - Installing EVE on Linux](https://wiki.eveuniversity.org/Installing_EVE_on_Linux)
- [ESI API Documentation](https://docs.esi.evetech.net/) - Public character endpoint for name resolution
