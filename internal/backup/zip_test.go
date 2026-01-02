package backup

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateAndReadBackup(t *testing.T) {
	tempDir := t.TempDir()

	// Create test source file
	sourceFile := filepath.Join(tempDir, "core_char_12345.dat")
	testContent := []byte("test character settings data")
	if err := os.WriteFile(sourceFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create backup
	backupPath := filepath.Join(tempDir, "test-backup.zip")
	chars := []CharacterBackup{
		{
			CharacterID:   12345,
			CharacterName: "Test Character",
			OriginalPath:  sourceFile,
			FileName:      "core_char_12345.dat",
		},
	}
	files := map[int64]string{
		12345: sourceFile,
	}

	if err := CreateBackup(backupPath, chars, files); err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Verify backup was created
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Fatal("backup file was not created")
	}

	// Read backup metadata
	metadata, err := ReadBackup(backupPath)
	if err != nil {
		t.Fatalf("ReadBackup failed: %v", err)
	}

	if metadata.Version != backupVersion {
		t.Errorf("expected version %s, got %s", backupVersion, metadata.Version)
	}

	if len(metadata.Characters) != 1 {
		t.Errorf("expected 1 character, got %d", len(metadata.Characters))
	}

	if metadata.Characters[0].CharacterName != "Test Character" {
		t.Errorf("expected 'Test Character', got '%s'", metadata.Characters[0].CharacterName)
	}
}

func TestExtractCharacter(t *testing.T) {
	tempDir := t.TempDir()

	// Create test source file
	sourceFile := filepath.Join(tempDir, "core_char_12345.dat")
	testContent := []byte("test character settings data")
	if err := os.WriteFile(sourceFile, testContent, 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create backup
	backupPath := filepath.Join(tempDir, "test-backup.zip")
	chars := []CharacterBackup{
		{
			CharacterID:   12345,
			CharacterName: "Test Character",
			OriginalPath:  sourceFile,
			FileName:      "core_char_12345.dat",
		},
	}
	files := map[int64]string{
		12345: sourceFile,
	}

	if err := CreateBackup(backupPath, chars, files); err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Extract character
	extractPath := filepath.Join(tempDir, "extracted", "core_char_12345.dat")
	if err := ExtractCharacter(backupPath, 12345, extractPath); err != nil {
		t.Fatalf("ExtractCharacter failed: %v", err)
	}

	// Verify extracted content
	extracted, err := os.ReadFile(extractPath)
	if err != nil {
		t.Fatalf("failed to read extracted file: %v", err)
	}

	if string(extracted) != string(testContent) {
		t.Errorf("extracted content mismatch: got %s, want %s", extracted, testContent)
	}
}

func TestExtractCharacterNotFound(t *testing.T) {
	tempDir := t.TempDir()

	// Create test source file
	sourceFile := filepath.Join(tempDir, "core_char_12345.dat")
	if err := os.WriteFile(sourceFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create backup
	backupPath := filepath.Join(tempDir, "test-backup.zip")
	chars := []CharacterBackup{
		{CharacterID: 12345, CharacterName: "Test", OriginalPath: sourceFile, FileName: "core_char_12345.dat"},
	}
	files := map[int64]string{12345: sourceFile}

	if err := CreateBackup(backupPath, chars, files); err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Try to extract non-existent character
	extractPath := filepath.Join(tempDir, "extracted", "core_char_99999.dat")
	err := ExtractCharacter(backupPath, 99999, extractPath)
	if err == nil {
		t.Error("expected error for non-existent character")
	}
}

func TestExtractAll(t *testing.T) {
	tempDir := t.TempDir()

	// Create test source files
	sourceFile1 := filepath.Join(tempDir, "core_char_111.dat")
	sourceFile2 := filepath.Join(tempDir, "core_char_222.dat")
	if err := os.WriteFile(sourceFile1, []byte("char1"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	if err := os.WriteFile(sourceFile2, []byte("char2"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Create backup
	backupPath := filepath.Join(tempDir, "test-backup.zip")
	chars := []CharacterBackup{
		{CharacterID: 111, CharacterName: "Char1", OriginalPath: sourceFile1, FileName: "core_char_111.dat"},
		{CharacterID: 222, CharacterName: "Char2", OriginalPath: sourceFile2, FileName: "core_char_222.dat"},
	}
	files := map[int64]string{
		111: sourceFile1,
		222: sourceFile2,
	}

	if err := CreateBackup(backupPath, chars, files); err != nil {
		t.Fatalf("CreateBackup failed: %v", err)
	}

	// Extract all
	extractDir := filepath.Join(tempDir, "extracted")
	if err := ExtractAll(backupPath, extractDir); err != nil {
		t.Fatalf("ExtractAll failed: %v", err)
	}

	// Verify both files extracted
	files1, _ := os.ReadFile(filepath.Join(extractDir, "core_char_111.dat"))
	files2, _ := os.ReadFile(filepath.Join(extractDir, "core_char_222.dat"))

	if string(files1) != "char1" {
		t.Errorf("char1 content mismatch")
	}
	if string(files2) != "char2" {
		t.Errorf("char2 content mismatch")
	}
}

func TestReadBackupInvalidFile(t *testing.T) {
	_, err := ReadBackup("/nonexistent/path.zip")
	if err == nil {
		t.Error("expected error for non-existent file")
	}
}
