package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// GetDataDir returns the platform-specific user data directory for MrRSS
func GetDataDir() (string, error) {
	var baseDir string
	var err error

	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("APPDATA")
		if baseDir == "" {
			baseDir = os.Getenv("USERPROFILE")
			if baseDir != "" {
				baseDir = filepath.Join(baseDir, "AppData", "Roaming")
			}
		}
	case "darwin":
		baseDir = os.Getenv("HOME")
		if baseDir != "" {
			baseDir = filepath.Join(baseDir, "Library", "Application Support")
		}
	case "linux":
		// Follow XDG Base Directory specification
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			homeDir := os.Getenv("HOME")
			if homeDir != "" {
				baseDir = filepath.Join(homeDir, ".local", "share")
			}
		}
	default:
		// Fallback for other platforms
		baseDir = os.Getenv("HOME")
		if baseDir != "" {
			baseDir = filepath.Join(baseDir, ".config")
		}
	}

	if baseDir == "" {
		// Last resort: use current directory
		baseDir, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	// Create MrRSS subdirectory
	dataDir := filepath.Join(baseDir, "MrRSS")
	err = os.MkdirAll(dataDir, 0755)
	if err != nil {
		return "", err
	}

	return dataDir, nil
}

// GetDBPath returns the full path to the database file
func GetDBPath() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "rss.db"), nil
}

// GetLogPath returns the full path to the debug log file
func GetLogPath() (string, error) {
	dataDir, err := GetDataDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dataDir, "debug.log"), nil
}
