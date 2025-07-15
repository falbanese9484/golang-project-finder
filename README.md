# Project Finder

A CLI tool to fuzzy find projects by directory name and open them in VSCode or tmux.

## Overview

Project Finder solves the common problem of managing multiple project directories scattered across your filesystem. Instead of navigating through nested folders or remembering exact project locations, you can quickly search and open projects using fuzzy matching.

The tool works by indexing directories in your chosen location (Desktop, Documents, or Downloads) and storing project metadata in a JSON file. When searching, it uses fuzzy matching to find projects and presents them in an interactive terminal interface, sorted by last modified date.

## Features

- **Fuzzy search**: Find projects by typing partial directory names
- **Interactive selection**: Browse through matching projects with arrow keys  
- **VSCode integration**: Open projects directly in VSCode
- **Workspace management**: Create and manage VS Code workspace files
- **Tmux support**: Navigate to project directory and run tmux-dev
- **Fast indexing**: Pre-index directories for quick searching
- **Smart filtering**: Excludes common directories like node_modules and venv

## Installation

1. **Clone the repository**:
   ```bash
   git clone https://github.com/falbanese9484/project-finder.git
   ```

2. **Navigate to the project directory**:
   ```bash
   cd project-finder
   ```

3. **Build the CLI tool**:
   ```bash
   go build -o findit
   ```

4. **Add the tool to your PATH** (optional):
   ```bash
   export PATH=$PATH:/path/to/project-finder
   ```

## Usage

### Initial Setup

1. **Configure the root directory**:
   ```bash
   ./findit config
   ```
   Choose from Desktop, Documents, or Downloads as your project root.

2. **Index your projects**:
   ```bash
   ./findit index
   ```
   This scans your configured directory and creates a searchable index.

### Finding Projects

**Open in VSCode** (default):
```bash
./findit find <search-query>
```

**Open in tmux**:
```bash
./findit find <search-query> --tmux
```

### Managing Workspaces

**Create a new workspace**:
```bash
./findit workspace init <workspace-name>
```

**Index existing workspaces**:
```bash
./findit workspace index
```

**Find and open a workspace**:
```bash
./findit workspace find <search-query>
```

## Commands

### Project Management
- `findit config` - Configure the root directory to index (Desktop, Documents, or Downloads)
- `findit index` - Index all directories in the configured location  
- `findit find <query>` - Search for projects matching the query and open in VSCode
- `findit find <query> --tmux` - Search for projects and open in tmux instead of VSCode

### Workspace Management
- `findit workspace init <name>` - Create a new VS Code workspace file
- `findit workspace index` - Index VS Code workspaces in ~/Desktop/workspaces
- `findit workspace find <query>` - Find and open a VS Code workspace

## Configuration

The tool creates a `.project-finder` directory in your home folder containing:
- `config.json` - Configuration settings (root directory selection)
- `projects.json` - Indexed project data with metadata
- `workspaces.json` - Indexed VS Code workspace files

Workspaces are stored in `~/Desktop/workspaces/` as `.code-workspace` files.

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [promptui](https://github.com/manifoldco/promptui) - Interactive terminal prompts
- [fuzzysearch](https://github.com/lithammer/fuzzysearch) - Fuzzy string matching

## Requirements

- Go 1.23.4+
- VSCode (for default mode)
- tmux and tmux-dev script (for tmux mode)

## Example Workflow

```bash
# Initial setup
./findit config          # Choose Desktop as root directory
./findit index           # Index all projects in ~/Desktop

# Daily usage  
./findit find react      # Find projects matching "react"
./findit find api --tmux # Find "api" projects and open in tmux

# Workspace management
./findit workspace init my-project    # Create new workspace
./findit workspace index              # Index existing workspaces
./findit workspace find my-project    # Find and open workspace
```
