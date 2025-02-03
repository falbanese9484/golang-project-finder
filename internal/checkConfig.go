package internal

import (
	"fmt"
	"os"
	"path/filepath"
)

func CheckConfig() error {
	rootDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(rootDir, ".project-finder", "config.json")
	if fileExists(configPath) {
		return nil
	} else {
		return fmt.Errorf("no config file found")
	}
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
