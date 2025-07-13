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

	"findit/internal"

	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func sortProjects(projects []Project) {
	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})
}

func readProjects() ([]Project, error) {
	rootDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(filepath.Join(rootDir, ".project-finder", "./projects.json"))
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(data, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func searchProjects(projects []Project, query string) []Project {
	var matches []Project

	for _, project := range projects {
		if fuzzy.Match(query, project.Name) {
			matches = append(matches, project)
		}
	}

	return matches
}

// findCmd represents the find command
var findCmd = &cobra.Command{
	Use:   "find [query]",
	Short: "Find a project (default: open new shell in directory, -c for VSCode, -t for tmux, -w for workspace)",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if the config file exists
		err := internal.CheckConfig()
		if err != nil {
			fmt.Println("Config file not found. Please run 'findit config' to set the config")
			return
		}

		workspaceQuery, _ := cmd.Flags().GetString("workspace")

		// Handle workspace search
		if workspaceQuery != "" {
			workspaces, err := ReadWorkspaces()
			if err != nil {
				fmt.Printf("Error reading workspaces: %v\n", err)
				fmt.Println("Try running 'findit workspace index' first")
				return
			}

			matches := SearchWorkspaces(workspaces, workspaceQuery)
			if len(matches) == 0 {
				fmt.Println("No workspaces found")
				return
			}

			SortWorkspaces(matches)

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
			return
		}

		// Handle project search (existing logic)
		if len(args) < 1 {
			fmt.Println("Please provide a search query")
			return
		}

		query := args[0]
		projects, err := readProjects()
		if err != nil {
			fmt.Printf("Error reading projects: %v\n", err)
			return
		}

		matches := searchProjects(projects, query)
		if len(matches) == 0 {
			fmt.Println("No projects found")
			return
		}

		sortProjects(matches)

		// Prepare a list of project names for interactive selection.
		var projectNames []string
		for _, p := range matches {
			if p.IsDir {
				projectNames = append(projectNames, fmt.Sprintf("%s -> Last_Modified: %s", p.Name, p.Modified.Format("2006-01-02 15:04:05")))
			}
		}

		// Using promptui to create an interactive terminal selector to cycle through found projects.
		selectPrompt := promptui.Select{
			Label: "Select a project",
			Items: projectNames,
		}
		index, _, err := selectPrompt.Run()
		if err != nil {
			fmt.Printf("Prompt failed: %v\n", err)
			return
		}

		selectedProject := matches[index]

		tmuxMode, _ := cmd.Flags().GetBool("tmux")
		codeMode, _ := cmd.Flags().GetBool("code")

		if tmuxMode {
			fmt.Printf("Navigating to project: %s at %s and running tmux-dev\n", selectedProject.Name, selectedProject.Path)

			// Change to the project directory and run tmux-dev
			tmuxCmd := exec.Command("bash", "-c", fmt.Sprintf("cd %s && tmux-dev", selectedProject.Path))
			tmuxCmd.Stdout = os.Stdout
			tmuxCmd.Stderr = os.Stderr
			tmuxCmd.Stdin = os.Stdin

			if err := tmuxCmd.Run(); err != nil {
				fmt.Printf("Error running tmux-dev: %v\n", err)
			}
		} else if codeMode {
			fmt.Printf("Opening project: %s at %s\n", selectedProject.Name, selectedProject.Path)

			// Attaches the selected directory to the current VS Code window.
			openCmd := exec.Command("code", "-a", selectedProject.Path)
			if err := openCmd.Start(); err != nil {
				fmt.Printf("Error opening project: %v\n", err)
			}
		} else {
			// Default behavior: change directory in current process and spawn a new shell
			fmt.Printf("Navigating to project: %s at %s\n", selectedProject.Name, selectedProject.Path)

			// Change the working directory of this process
			if err := os.Chdir(selectedProject.Path); err != nil {
				fmt.Printf("Error changing directory: %v\n", err)
				return
			}

			// Spawn a new shell in the target directory
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = "/bin/bash" // fallback
			}

			shellCmd := exec.Command(shell)
			shellCmd.Stdout = os.Stdout
			shellCmd.Stderr = os.Stderr
			shellCmd.Stdin = os.Stdin
			shellCmd.Dir = selectedProject.Path

			if err := shellCmd.Run(); err != nil {
				fmt.Printf("Error starting shell: %v\n", err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(findCmd)
	findCmd.Flags().BoolP("tmux", "t", false, "Navigate to directory and run tmux-dev")
	findCmd.Flags().BoolP("code", "c", false, "Open project in VSCode")
	findCmd.Flags().StringP("workspace", "w", "", "Search for VS Code workspace instead of project")
}
