package eve

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// CopySettings copies settings from one character to another.
// It creates a backup of the target file before overwriting.
func CopySettings(from, to *CharacterSettings, backupDir string) error {
	// Create backup of target if it exists
	if _, err := os.Stat(to.FilePath); err == nil {
		if backupDir != "" {
			backupPath := filepath.Join(backupDir, fmt.Sprintf("core_char_%d_%s.dat.bak",
				to.CharacterID, time.Now().Format("20060102_150405")))
			if err := copyFile(to.FilePath, backupPath); err != nil {
				return fmt.Errorf("failed to create backup: %w", err)
			}
		}
	}

	// Copy source to target
	return copyFile(from.FilePath, to.FilePath)
}

// GetSettingsDir returns the directory containing the settings file.
func (cs *CharacterSettings) GetSettingsDir() string {
	return filepath.Dir(cs.FilePath)
}

// copyFile copies a file from src to dst.
func copyFile(src, dst string) error {
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return err
	}

	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	return err
}

// CreateCharacterSettingsPath generates a path for a new character settings file.
// Uses the same directory as the reference character.
func CreateCharacterSettingsPath(referenceChar *CharacterSettings, newCharID int64) string {
	dir := referenceChar.GetSettingsDir()
	return filepath.Join(dir, fmt.Sprintf("core_char_%d.dat", newCharID))
}
