package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jpbriend/eve-settings-manager/internal/backup"
	"github.com/jpbriend/eve-settings-manager/internal/esi"
	"github.com/jpbriend/eve-settings-manager/internal/eve"
	"github.com/spf13/cobra"
)

var (
	copyFrom  string
	copyTo    string
	copyForce bool
)

var copyCmd = &cobra.Command{
	Use:   "copy",
	Short: "Copy settings from one character to another",
	Long: `Copy character settings from one character to another.

Works across different accounts. Automatically creates a backup of the target
character settings before overwriting.`,
	RunE: runCopy,
}

func init() {
	copyCmd.Flags().StringVar(&copyFrom, "from", "", "Source character (ID or name)")
	copyCmd.Flags().StringVar(&copyTo, "to", "", "Target character (ID or name)")
	copyCmd.Flags().BoolVarP(&copyForce, "force", "f", false, "Overwrite without confirmation")
	_ = copyCmd.MarkFlagRequired("from")
	_ = copyCmd.MarkFlagRequired("to")
}

func runCopy(cmd *cobra.Command, args []string) error {
	// ESI client for name resolution
	esiClient := esi.NewClient()

	// Resolve character IDs (supports both ID and name)
	fromID, err := esiClient.ResolveCharacter(copyFrom)
	if err != nil {
		return fmt.Errorf("failed to resolve source character '%s': %w", copyFrom, err)
	}

	toID, err := esiClient.ResolveCharacter(copyTo)
	if err != nil {
		return fmt.Errorf("failed to resolve target character '%s': %w", copyTo, err)
	}

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

	// Find source character
	var sourceChar *eve.CharacterSettings
	for _, c := range allCharacters {
		if c.CharacterID == fromID {
			sourceChar = &c
			break
		}
	}

	if sourceChar == nil {
		return fmt.Errorf("source character %d not found in local settings", fromID)
	}

	// Find or prepare target character
	var targetChar *eve.CharacterSettings
	for _, c := range allCharacters {
		if c.CharacterID == toID {
			targetChar = &c
			break
		}
	}

	// Get character names for display
	sourceName := esiClient.GetCharacterNameOrFallback(fromID)
	targetName := esiClient.GetCharacterNameOrFallback(toID)

	// If target doesn't exist locally, we need to create it
	var targetPath string
	if targetChar == nil {
		// Use same settings directory as source
		targetPath = eve.CreateCharacterSettingsPath(sourceChar, toID)
		fmt.Printf("Target character settings file will be created at:\n  %s\n", targetPath)
	} else {
		targetPath = targetChar.FilePath
	}

	// Confirmation prompt
	if !copyForce {
		fmt.Printf("\nAbout to copy settings:\n")
		fmt.Printf("  From: %s (%d)\n", sourceName, fromID)
		fmt.Printf("  To:   %s (%d)\n", targetName, toID)
		if targetChar != nil {
			fmt.Printf("\nWARNING: This will overwrite existing settings for %s\n", targetName)
		}
		fmt.Print("\nProceed? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	// Create backup of target if it exists
	if targetChar != nil {
		backupDir := filepath.Dir(targetChar.FilePath)
		backupPath := filepath.Join(backupDir, fmt.Sprintf("core_char_%d_%s.dat.bak",
			toID, time.Now().Format("20060102_150405")))

		charBackup := []backup.CharacterBackup{{
			CharacterID:   toID,
			CharacterName: targetName,
			OriginalPath:  targetChar.FilePath,
			FileName:      fmt.Sprintf("core_char_%d.dat", toID),
		}}
		files := map[int64]string{toID: targetChar.FilePath}

		zipBackupPath := filepath.Join(backupDir, fmt.Sprintf("backup_%d_%s.zip",
			toID, time.Now().Format("20060102_150405")))

		if err := backup.CreateBackup(zipBackupPath, charBackup, files); err != nil {
			return fmt.Errorf("failed to create backup of target: %w", err)
		}
		fmt.Printf("Backup created: %s\n", backupPath)
	}

	// Perform the copy
	targetSettings := &eve.CharacterSettings{
		CharacterID: toID,
		FilePath:    targetPath,
	}

	if err := eve.CopySettings(sourceChar, targetSettings, ""); err != nil {
		return fmt.Errorf("failed to copy settings: %w", err)
	}

	fmt.Printf("\nSettings copied successfully!\n")
	fmt.Printf("  From: %s (%d)\n", sourceName, fromID)
	fmt.Printf("  To:   %s (%d)\n", targetName, toID)

	return nil
}
