package backup

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Metadata contains backup metadata.
type Metadata struct {
	CreatedAt  string            `json:"created_at"`
	Version    string            `json:"version"`
	Characters []CharacterBackup `json:"characters"`
}

// CharacterBackup contains information about a backed up character.
type CharacterBackup struct {
	CharacterID   int64  `json:"character_id"`
	CharacterName string `json:"character_name"`
	OriginalPath  string `json:"original_path"`
	FileName      string `json:"file_name"`
}

const metadataFileName = "metadata.json"
const backupVersion = "1.0"

// CreateBackup creates a ZIP backup containing the specified character files.
func CreateBackup(outputPath string, characters []CharacterBackup, files map[int64]string) (err error) {
	zipFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer func() {
		if cerr := zipFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	zipWriter := zip.NewWriter(zipFile)
	defer func() {
		if cerr := zipWriter.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	// Create metadata
	metadata := Metadata{
		CreatedAt:  time.Now().Format(time.RFC3339),
		Version:    backupVersion,
		Characters: characters,
	}

	// Write metadata
	metadataWriter, err := zipWriter.Create(metadataFileName)
	if err != nil {
		return fmt.Errorf("failed to create metadata entry: %w", err)
	}
	if err := json.NewEncoder(metadataWriter).Encode(metadata); err != nil {
		return fmt.Errorf("failed to write metadata: %w", err)
	}

	// Write character files
	for charID, filePath := range files {
		fileName := fmt.Sprintf("core_char_%d.dat", charID)
		if err := addFileToZip(zipWriter, filePath, fileName); err != nil {
			return fmt.Errorf("failed to add character %d to backup: %w", charID, err)
		}
	}

	return nil
}

// ReadBackup reads and validates a backup file, returning its metadata.
func ReadBackup(backupPath string) (*Metadata, error) {
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() {
		_ = zipReader.Close()
	}()

	// Find and read metadata
	for _, file := range zipReader.File {
		if file.Name == metadataFileName {
			rc, err := file.Open()
			if err != nil {
				return nil, fmt.Errorf("failed to read metadata: %w", err)
			}

			var metadata Metadata
			decodeErr := json.NewDecoder(rc).Decode(&metadata)
			closeErr := rc.Close()

			if decodeErr != nil {
				return nil, fmt.Errorf("failed to parse metadata: %w", decodeErr)
			}
			if closeErr != nil {
				return nil, fmt.Errorf("failed to close metadata reader: %w", closeErr)
			}
			return &metadata, nil
		}
	}

	return nil, fmt.Errorf("backup file is missing metadata")
}

// ExtractCharacter extracts a specific character's settings from a backup.
func ExtractCharacter(backupPath string, charID int64, destPath string) error {
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() {
		_ = zipReader.Close()
	}()

	fileName := fmt.Sprintf("core_char_%d.dat", charID)
	for _, file := range zipReader.File {
		if file.Name == fileName {
			return extractFile(file, destPath)
		}
	}

	return fmt.Errorf("character %d not found in backup", charID)
}

// ExtractAll extracts all character settings from a backup to a directory.
func ExtractAll(backupPath, destDir string) error {
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup file: %w", err)
	}
	defer func() {
		_ = zipReader.Close()
	}()

	for _, file := range zipReader.File {
		if file.Name == metadataFileName {
			continue
		}
		destPath := filepath.Join(destDir, file.Name)
		if err := extractFile(file, destPath); err != nil {
			return err
		}
	}

	return nil
}

func addFileToZip(zipWriter *zip.Writer, srcPath, destName string) (err error) {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := srcFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	info, err := srcFile.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}
	header.Name = destName
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}

	_, err = io.Copy(writer, srcFile)
	return err
}

func extractFile(file *zip.File, destPath string) (err error) {
	if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
		return err
	}

	rc, err := file.Open()
	if err != nil {
		return err
	}
	defer func() {
		if cerr := rc.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := destFile.Close(); cerr != nil && err == nil {
			err = cerr
		}
	}()

	_, err = io.Copy(destFile, rc)
	return err
}
