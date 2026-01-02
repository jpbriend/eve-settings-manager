package commands

import (
	"fmt"
	"time"

	"github.com/jpbriend/eve-settings-manager/internal/backup"
	"github.com/jpbriend/eve-settings-manager/internal/esi"
	"github.com/jpbriend/eve-settings-manager/internal/eve"
	"github.com/spf13/cobra"
)

var (
	backupAll    bool
	backupOutput string
)

var backupCmd = &cobra.Command{
	Use:   "backup [character]",
	Short: "Backup character settings to a ZIP file",
	Long: `Create a ZIP backup of character settings.

You can specify a character by ID or name. Use --all to backup all characters.
The backup includes metadata with character names and timestamps.`,
	RunE: runBackup,
}

func init() {
	backupCmd.Flags().BoolVar(&backupAll, "all", false, "Backup all characters")
	backupCmd.Flags().StringVarP(&backupOutput, "output", "o", "", "Output file path")
}

func runBackup(cmd *cobra.Command, args []string) error {
	// Detect settings directories
	dirs, err := eve.DetectSettingsDirectories()
	if err != nil {
		return fmt.Errorf("failed to detect settings directories: %w", err)
	}

	if len(dirs) == 0 {
		return fmt.Errorf("no Eve Online settings directories found")
	}

	// Find all character settings
	allCharacters, err := eve.FindCharacterSettings(dirs)
	if err != nil {
		return fmt.Errorf("failed to find character settings: %w", err)
	}

	if len(allCharacters) == 0 {
		return fmt.Errorf("no character settings files found")
	}

	// ESI client for name resolution
	esiClient := esi.NewClient()

	// Determine which characters to backup
	var charactersToBackup []eve.CharacterSettings

	if backupAll {
		charactersToBackup = allCharacters
	} else if len(args) > 0 {
		// Resolve character by ID or name
		charID, err := esiClient.ResolveCharacter(args[0])
		if err != nil {
			return fmt.Errorf("failed to resolve character '%s': %w", args[0], err)
		}

		for _, c := range allCharacters {
			if c.CharacterID == charID {
				charactersToBackup = append(charactersToBackup, c)
				break
			}
		}

		if len(charactersToBackup) == 0 {
			return fmt.Errorf("character '%s' (ID: %d) not found in local settings", args[0], charID)
		}
	} else {
		return fmt.Errorf("please specify a character (ID or name) or use --all to backup all characters")
	}

	// Fetch character names
	charIDs := make([]int64, len(charactersToBackup))
	for i, c := range charactersToBackup {
		charIDs[i] = c.CharacterID
	}
	names := esiClient.BatchGetCharacterNames(charIDs)

	// Prepare backup data
	backupChars := make([]backup.CharacterBackup, len(charactersToBackup))
	files := make(map[int64]string)

	for i, c := range charactersToBackup {
		backupChars[i] = backup.CharacterBackup{
			CharacterID:   c.CharacterID,
			CharacterName: names[c.CharacterID],
			OriginalPath:  c.FilePath,
			FileName:      fmt.Sprintf("core_char_%d.dat", c.CharacterID),
		}
		files[c.CharacterID] = c.FilePath
	}

	// Determine output path
	outputPath := backupOutput
	if outputPath == "" {
		outputPath = fmt.Sprintf("eve-backup-%s.zip", time.Now().Format("20060102-150405"))
	}

	// Create backup
	if err := backup.CreateBackup(outputPath, backupChars, files); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	fmt.Printf("Backup created: %s\n", outputPath)
	fmt.Printf("Characters backed up: %d\n", len(charactersToBackup))
	for _, c := range backupChars {
		fmt.Printf("  - %s (%d)\n", c.CharacterName, c.CharacterID)
	}

	return nil
}
