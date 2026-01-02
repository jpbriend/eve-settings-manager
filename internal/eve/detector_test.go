package eve

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindCharacterSettings(t *testing.T) {
	// Create temp directory structure
	tempDir := t.TempDir()
	settingsDir := filepath.Join(tempDir, "settings_Default")
	if err := os.MkdirAll(settingsDir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}

	// Create test character files
	testFiles := []struct {
		name     string
		charID   int64
		isChar   bool
	}{
		{"core_char_12345678.dat", 12345678, true},
		{"core_char_87654321.dat", 87654321, true},
		{"core_user_99999.dat", 0, false},
		{"other_file.txt", 0, false},
	}

	for _, tf := range testFiles {
		path := filepath.Join(settingsDir, tf.name)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	// Test FindCharacterSettings
	chars, err := FindCharacterSettings([]string{settingsDir})
	if err != nil {
		t.Fatalf("FindCharacterSettings failed: %v", err)
	}

	if len(chars) != 2 {
		t.Errorf("expected 2 characters, got %d", len(chars))
	}

	// Verify character IDs
	foundIDs := make(map[int64]bool)
	for _, c := range chars {
		foundIDs[c.CharacterID] = true
	}

	for _, tf := range testFiles {
		if tf.isChar {
			if !foundIDs[tf.charID] {
				t.Errorf("expected to find character %d", tf.charID)
			}
		}
	}
}

func TestDetectSettingsDirectories(t *testing.T) {
	// This test verifies the function doesn't panic on various systems
	// Actual directories may or may not exist
	dirs, err := DetectSettingsDirectories()
	if err != nil {
		t.Errorf("DetectSettingsDirectories returned error: %v", err)
	}

	// Result can be empty (no Eve installed) or contain paths
	t.Logf("Found %d settings directories", len(dirs))
	for _, dir := range dirs {
		t.Logf("  - %s", dir)
	}
}

func TestCharacterSettingsGetSettingsDir(t *testing.T) {
	cs := &CharacterSettings{
		CharacterID: 12345,
		FilePath:    "/path/to/settings_Default/core_char_12345.dat",
	}

	expected := "/path/to/settings_Default"
	if got := cs.GetSettingsDir(); got != expected {
		t.Errorf("GetSettingsDir() = %s, want %s", got, expected)
	}
}

func TestCreateCharacterSettingsPath(t *testing.T) {
	ref := &CharacterSettings{
		CharacterID: 12345,
		FilePath:    "/path/to/settings_Default/core_char_12345.dat",
	}

	newPath := CreateCharacterSettingsPath(ref, 67890)
	expected := "/path/to/settings_Default/core_char_67890.dat"

	if newPath != expected {
		t.Errorf("CreateCharacterSettingsPath() = %s, want %s", newPath, expected)
	}
}
