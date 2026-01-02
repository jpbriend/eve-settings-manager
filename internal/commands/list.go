package commands

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/jpbriend/eve-settings-manager/internal/esi"
	"github.com/jpbriend/eve-settings-manager/internal/eve"
	"github.com/spf13/cobra"
)

var listVerbose bool

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all detected Eve character settings",
	Long: `List all detected Eve Online character settings files.

Scans known Eve settings locations and displays character IDs with their names
(resolved via ESI API), modification times, and file paths.`,
	RunE: runList,
}

func init() {
	listCmd.Flags().BoolVarP(&listVerbose, "verbose", "v", false, "Show additional details including full paths")
}

func runList(cmd *cobra.Command, args []string) error {
	// Detect settings directories
	dirs, err := eve.DetectSettingsDirectories()
	if err != nil {
		return fmt.Errorf("failed to detect settings directories: %w", err)
	}

	if len(dirs) == 0 {
		fmt.Println("No Eve Online settings directories found.")
		fmt.Println("\nSearched locations:")
		for _, path := range eve.GetPossibleSettingsPaths() {
			fmt.Printf("  - %s\n", path)
		}
		return nil
	}

	if listVerbose {
		fmt.Println("Found settings directories:")
		for _, dir := range dirs {
			fmt.Printf("  - %s\n", dir)
		}
		fmt.Println()
	}

	// Find character settings files
	characters, err := eve.FindCharacterSettings(dirs)
	if err != nil {
		return fmt.Errorf("failed to find character settings: %w", err)
	}

	if len(characters) == 0 {
		fmt.Println("No character settings files found.")
		return nil
	}

	// Fetch character names from ESI
	client := esi.NewClient()
	charIDs := make([]int64, len(characters))
	for i, c := range characters {
		charIDs[i] = c.CharacterID
	}
	names := client.BatchGetCharacterNames(charIDs)

	// Display results
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	if listVerbose {
		fmt.Fprintln(w, "CHARACTER ID\tNAME\tMODIFIED\tPATH")
	} else {
		fmt.Fprintln(w, "CHARACTER ID\tNAME\tMODIFIED")
	}

	for _, c := range characters {
		name := names[c.CharacterID]
		modTime := time.Unix(c.ModTime, 0).Format("2006-01-02 15:04:05")

		if listVerbose {
			fmt.Fprintf(w, "%d\t%s\t%s\t%s\n", c.CharacterID, name, modTime, c.FilePath)
		} else {
			fmt.Fprintf(w, "%d\t%s\t%s\n", c.CharacterID, name, modTime)
		}
	}
	w.Flush()

	fmt.Printf("\nFound %d character(s)\n", len(characters))
	return nil
}
