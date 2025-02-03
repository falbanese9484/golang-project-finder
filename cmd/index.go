/*
Copyright © 2025 Frank Albanese <albanesefc9@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"project-finder/internal"

	"github.com/spf13/cobra"
)

type Project struct {
	Name     string
	Path     string
	IsDir    bool
	Modified time.Time
}

func getProjects(dir string) ([]Project, error) {
	var projects []Project
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if strings.Contains(path, "node_modules") || strings.Contains(path, "venv") {
			return nil
		}
		if info.IsDir() {
			project := Project{
				Name:     info.Name(),
				Path:     path,
				IsDir:    info.IsDir(),
				Modified: info.ModTime(),
			}
			projects = append(projects, project)
		}
		return nil
	})
	return projects, err
}

func writeProjectsToFile(projects []Project, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(projects)
	if err != nil {
		return err
	}

	return nil
}

func getCurrentTime() time.Time {
	return time.Now()
}

// indexCmd represents the index command
var indexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index files in preset Directory - ~/Desktop/Projects",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.CheckConfig()
		if err != nil {
			fmt.Println("No config file found. Please run 'findit config' to create one.")
			return
		}
		startTime := getCurrentTime()
		rootPath, err := os.UserHomeDir()
		projectsPath := filepath.Join(rootPath, "Desktop")
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}
		projects, err := getProjects(projectsPath)
		if err != nil {
			fmt.Println("Error getting projects:", err)
			return
		}
		err = writeProjectsToFile(projects, filepath.Join(rootPath, ".project-finder", "projects.json"))
		if err != nil {
			fmt.Println("Error writing projects to file:", err)
			return
		}
		endTime := getCurrentTime()
		fmt.Println("Indexing completed in", endTime.Sub(startTime))
	},
}

func init() {
	rootCmd.AddCommand(indexCmd)
}
