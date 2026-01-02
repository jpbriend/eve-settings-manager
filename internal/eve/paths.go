package eve

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetPossibleSettingsPaths returns all possible Eve settings base paths for the current platform.
func GetPossibleSettingsPaths() []string {
	switch runtime.GOOS {
	case "windows":
		return getWindowsPaths()
	case "linux":
		return getLinuxPaths()
	default:
		return nil
	}
}

func getWindowsPaths() []string {
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		return nil
	}
	return []string{
		filepath.Join(localAppData, "CCP", "EVE"),
	}
}

func getLinuxPaths() []string {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil
	}

	paths := []string{}

	// Steam Proton paths
	steamPaths := []string{
		filepath.Join(home, ".steam", "steam", "steamapps", "compatdata", "8500", "pfx", "drive_c", "users", "steamuser", "AppData", "Local", "CCP", "EVE"),
		filepath.Join(home, ".local", "share", "Steam", "steamapps", "compatdata", "8500", "pfx", "drive_c", "users", "steamuser", "AppData", "Local", "CCP", "EVE"),
	}
	paths = append(paths, steamPaths...)

	// Lutris paths - check common locations
	lutrisPaths := []string{
		filepath.Join(home, "Games", "eve-online", "drive_c", "users", home, "AppData", "Local", "CCP", "EVE"),
		filepath.Join(home, ".wine", "drive_c", "users", home, "AppData", "Local", "CCP", "EVE"),
	}
	paths = append(paths, lutrisPaths...)

	return paths
}
