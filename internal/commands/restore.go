package commands

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jpbriend/eve-settings-manager/internal/backup"
	"github.com/jpbriend/eve-settings-manager/internal/eve"
	"github.com/spf13/cobra"
)

var (
	restoreCharacter string
	restoreForce     bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore <backup.zip>",
	Short: "Restore character settings from a backup",
	Long: `Restore character settings from a ZIP backup file.

By default, restores all characters in the backup to their original locations.
Use --character to restore a specific character only.`,
	Args: cobra.ExactArgs(1),
	RunE: runRestore,
}

func init() {
	restoreCmd.Flags().StringVarP(&restoreCharacter, "character", "c", "", "Restore specific character (ID or name)")
	restoreCmd.Flags().BoolVarP(&restoreForce, "force", "f", false, "Restore without confirmation")
}

func runRestore(cmd *cobra.Command, args []string) error {
	backupFile := args[0]

	// Read backup metadata
	metadata, err := backup.ReadBackup(backupFile)
	if err != nil {
		return fmt.Errorf("failed to read backup: %w", err)
	}

	fmt.Printf("Backup file: %s\n", backupFile)
	fmt.Printf("Created: %s\n", metadata.CreatedAt)
	fmt.Printf("Version: %s\n", metadata.Version)
	fmt.Printf("Characters in backup:\n")
	for _, c := range metadata.Characters {
		fmt.Printf("  - %s (%d)\n", c.CharacterName, c.CharacterID)
	}

	// Determine which characters to restore
	var charactersToRestore []backup.CharacterBackup

	if restoreCharacter != "" {
		// Find specific character
		charID, err := strconv.ParseInt(restoreCharacter, 10, 64)
		if err != nil {
			// Try to find by name
			for _, c := range metadata.Characters {
				if strings.EqualFold(c.CharacterName, restoreCharacter) {
					charactersToRestore = append(charactersToRestore, c)
					break
				}
			}
		} else {
			for _, c := range metadata.Characters {
				if c.CharacterID == charID {
					charactersToRestore = append(charactersToRestore, c)
					break
				}
			}
		}

		if len(charactersToRestore) == 0 {
			return fmt.Errorf("character '%s' not found in backup", restoreCharacter)
		}
	} else {
		charactersToRestore = metadata.Characters
	}

	// Check if we have Eve settings directories to restore to
	dirs, err := eve.DetectSettingsDirectories()
	if err != nil {
		return fmt.Errorf("failed to detect settings directories: %w", err)
	}

	if len(dirs) == 0 {
		return fmt.Errorf("no Eve Online settings directories found - cannot restore")
	}

	// Determine restore paths
	fmt.Printf("\nWill restore to:\n")
	restorePaths := make(map[int64]string)
	for _, c := range charactersToRestore {
		// Try to use original path if it exists and is in a valid settings dir
		restorePath := c.OriginalPath
		pathValid := false

		for _, dir := range dirs {
			if strings.HasPrefix(c.OriginalPath, dir) {
				pathValid = true
				break
			}
		}

		if !pathValid {
			// Use first available settings directory
			restorePath = fmt.Sprintf("%s/core_char_%d.dat", dirs[0], c.CharacterID)
		}

		restorePaths[c.CharacterID] = restorePath
		fmt.Printf("  %s (%d) -> %s\n", c.CharacterName, c.CharacterID, restorePath)
	}

	// Confirmation prompt
	if !restoreForce {
		fmt.Print("\nProceed with restore? [y/N]: ")

		reader := bufio.NewReader(os.Stdin)
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("Operation cancelled.")
			return nil
		}
	}

	// Perform restore
	for _, c := range charactersToRestore {
		destPath := restorePaths[c.CharacterID]
		if err := backup.ExtractCharacter(backupFile, c.CharacterID, destPath); err != nil {
			return fmt.Errorf("failed to restore character %d: %w", c.CharacterID, err)
		}
		fmt.Printf("Restored: %s (%d)\n", c.CharacterName, c.CharacterID)
	}

	fmt.Printf("\nRestore completed successfully! %d character(s) restored.\n", len(charactersToRestore))
	return nil
}
