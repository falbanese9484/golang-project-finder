/*
Copyright © 2025 Frank Albanese <albanesefc9@gmail.com>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"findit/internal"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Workspace struct {
	Name     string    `json:"name"`
	Path     string    `json:"path"`
	Modified time.Time `json:"modified"`
}

func getWorkspaces(dir string) ([]Workspace, error) {
	var workspaces []Workspace

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".code-workspace") {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			workspace := Workspace{
				Name:     strings.TrimSuffix(entry.Name(), ".code-workspace"),
				Path:     filepath.Join(dir, entry.Name()),
				Modified: info.ModTime(),
			}
			workspaces = append(workspaces, workspace)
		}
	}

	return workspaces, nil
}

func writeWorkspacesToFile(workspaces []Workspace, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(workspaces)
	if err != nil {
		return err
	}

	return nil
}

func readWorkspaces() ([]Workspace, error) {
	rootDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(filepath.Join(rootDir, ".project-finder", "workspaces.json"))
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	err = json.Unmarshal(data, &workspaces)
	if err != nil {
		return nil, err
	}

	return workspaces, nil
}

func searchWorkspaces(workspaces []Workspace, query string) []Workspace {
	var matches []Workspace

	for _, workspace := range workspaces {
		if fuzzy.Match(query, workspace.Name) {
			matches = append(matches, workspace)
		}
	}

	return matches
}

func sortWorkspaces(workspaces []Workspace) {
	sort.Slice(workspaces, func(i, j int) bool {
		return workspaces[i].Modified.After(workspaces[j].Modified)
	})
}

// workspaceCmd represents the workspace command
var workspaceCmd = &cobra.Command{
	Use:   "workspace",
	Short: "Manage VS Code workspaces",
}

// workspaceIndexCmd represents the workspace index command
var workspaceIndexCmd = &cobra.Command{
	Use:   "index",
	Short: "Index VS Code workspaces in ~/Desktop/workspaces",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.CheckConfig()
		if err != nil {
			fmt.Println("Config file not found. Please run 'findit config' to set the config")
			return
		}

		startTime := time.Now()
		rootDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}

		workspacesPath := filepath.Join(rootDir, "Desktop", "workspaces")
		if _, err := os.Stat(workspacesPath); os.IsNotExist(err) {
			fmt.Printf("Workspaces directory not found: %s\n", workspacesPath)
			return
		}

		workspaces, err := getWorkspaces(workspacesPath)
		if err != nil {
			fmt.Println("Error getting workspaces:", err)
			return
		}

		configDir := filepath.Join(rootDir, ".project-finder")
		if _, err := os.Stat(configDir); os.IsNotExist(err) {
			err = os.MkdirAll(configDir, 0755)
			if err != nil {
				fmt.Println("Error creating config directory:", err)
				return
			}
		}

		err = writeWorkspacesToFile(workspaces, filepath.Join(configDir, "workspaces.json"))
		if err != nil {
			fmt.Println("Error writing workspaces to file:", err)
			return
		}

		endTime := time.Now()
		fmt.Printf("Indexed %d workspaces in %v\n", len(workspaces), endTime.Sub(startTime))
	},
}

// workspaceFindCmd represents the workspace find command
var workspaceFindCmd = &cobra.Command{
	Use:   "find [query]",
	Short: "Find and open a VS Code workspace",
	Run: func(cmd *cobra.Command, args []string) {
		err := internal.CheckConfig()
		if err != nil {
			fmt.Println("Config file not found. Please run 'findit config' to set the config")
			return
		}

		if len(args) < 1 {
			fmt.Println("Please provide a search query")
			return
		}

		query := args[0]
		workspaces, err := readWorkspaces()
		if err != nil {
			fmt.Printf("Error reading workspaces: %v\n", err)
			fmt.Println("Try running 'findit workspace index' first")
			return
		}

		matches := searchWorkspaces(workspaces, query)
		if len(matches) == 0 {
			fmt.Println("No workspaces found")
			return
		}

		sortWorkspaces(matches)

		var workspaceNames []string
		for _, w := range matches {
			workspaceNames = append(workspaceNames, fmt.Sprintf("%s -> Last_Modified: %s", w.Name, w.Modified.Format("2006-01-02 15:04:05")))
		}

		selectPrompt := promptui.Select{
			Label: "Select a workspace",
			Items: workspaceNames,
		}

		index, _, err := selectPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		selectedWorkspace := matches[index]
		fmt.Printf("Opening workspace: %s\n", selectedWorkspace.Name)

		openCmd := exec.Command("code", selectedWorkspace.Path)
		if err := openCmd.Start(); err != nil {
			fmt.Printf("Error opening workspace: %v\n", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(workspaceCmd)
	workspaceCmd.AddCommand(workspaceIndexCmd)
	workspaceCmd.AddCommand(workspaceFindCmd)
}
