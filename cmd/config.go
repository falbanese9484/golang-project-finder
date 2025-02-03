/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type Config struct {
	RootDir string `json:"rootDir"`
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set the config for your finder.",
	Run: func(cmd *cobra.Command, args []string) {
		rootDir, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory")
			return
		}

		configPath := rootDir + "/.project-finder/config.json"
		if fileExists(configPath) {
			fmt.Println("Config already exists")
			return
		}

		selectPrompt := promptui.Select{
			Label: "Select the directory to index",
			Items: []string{"Desktop", "Documents", "Downloads"},
		}

		_, result, err := selectPrompt.Run()
		if err != nil {
			fmt.Println("Error selecting directory")
			return
		}

		config := Config{
			RootDir: result,
		}

		err = writeConfigToFile(config, configPath)

		if err != nil {
			fmt.Println("Error writing config to file")
			return
		}

	},
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func writeConfigToFile(config Config, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
