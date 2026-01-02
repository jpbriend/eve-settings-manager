package eve

import (
	"os"
	"path/filepath"
	"regexp"
	"strconv"
)

// CharacterSettings represents a character's settings file.
type CharacterSettings struct {
	CharacterID int64
	FilePath    string
	ModTime     int64 // Unix timestamp
}

// DetectSettingsDirectories finds all Eve settings directories.
func DetectSettingsDirectories() ([]string, error) {
	var settingsDirs []string
	basePaths := GetPossibleSettingsPaths()

	for _, basePath := range basePaths {
		if _, err := os.Stat(basePath); os.IsNotExist(err) {
			continue
		}

		// Look for directories containing settings_* subdirectories
		entries, err := os.ReadDir(basePath)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			// Check for settings_* subdirectories (e.g., settings_Default)
			profilePath := filepath.Join(basePath, entry.Name())
			profileEntries, err := os.ReadDir(profilePath)
			if err != nil {
				continue
			}

			for _, profileEntry := range profileEntries {
				if profileEntry.IsDir() && regexp.MustCompile(`^settings_`).MatchString(profileEntry.Name()) {
					settingsDir := filepath.Join(profilePath, profileEntry.Name())
					settingsDirs = append(settingsDirs, settingsDir)
				}
			}
		}
	}

	return settingsDirs, nil
}

// FindCharacterSettings finds all core_char_*.dat files in the given directories.
func FindCharacterSettings(settingsDirs []string) ([]CharacterSettings, error) {
	var characters []CharacterSettings
	charFilePattern := regexp.MustCompile(`^core_char_(\d+)\.dat$`)

	for _, dir := range settingsDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			matches := charFilePattern.FindStringSubmatch(entry.Name())
			if matches == nil {
				continue
			}

			charID, err := strconv.ParseInt(matches[1], 10, 64)
			if err != nil {
				continue
			}

			info, err := entry.Info()
			if err != nil {
				continue
			}

			characters = append(characters, CharacterSettings{
				CharacterID: charID,
				FilePath:    filepath.Join(dir, entry.Name()),
				ModTime:     info.ModTime().Unix(),
			})
		}
	}

	return characters, nil
}

// FindCharacterByID finds a character settings file by character ID.
func FindCharacterByID(charID int64) (*CharacterSettings, error) {
	dirs, err := DetectSettingsDirectories()
	if err != nil {
		return nil, err
	}

	characters, err := FindCharacterSettings(dirs)
	if err != nil {
		return nil, err
	}

	for _, char := range characters {
		if char.CharacterID == charID {
			return &char, nil
		}
	}

	return nil, nil
}
