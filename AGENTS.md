# AGENTS.md - Development Guidelines

## Build/Test Commands
- **Build**: `go build -o findit .` or use `./build.sh` for installation
- **Run**: `./findit <command>` or `findit <command>` if in PATH
- **Test**: `go test ./...` (no tests currently exist)
- **Single test**: `go test -run TestFunctionName ./package`
- **Vet**: `go vet ./...`

## Available Commands
### Project Management
- `findit config` - Configure root directory (Desktop/Documents/Downloads)
- `findit index` - Index directories in configured location
- `findit find <query>` - Find and open projects in VSCode
- `findit find <query> --tmux` - Find and open projects in tmux

### Workspace Management
- `findit workspace init <name>` - Create new VS Code workspace file
- `findit workspace index` - Index existing workspace files in ~/Desktop/workspaces
- `findit workspace find <query>` - Find and open VS Code workspace

## Code Style Guidelines
- **Module**: `findit` (as defined in go.mod)
- **Go version**: 1.23.4+
- **Imports**: Standard library first, then third-party, then local packages
- **Package structure**: `cmd/` for CLI commands, `internal/` for internal logic
- **Types**: Use struct tags for JSON serialization (`json:"fieldName"`)
- **Error handling**: Return errors, check with `if err != nil`, print user-friendly messages
- **Naming**: CamelCase for exported, camelCase for unexported, descriptive names
- **Comments**: Use copyright headers, document exported functions/types
- **File operations**: Use `filepath.Join()` for paths, check file existence with `os.Stat()`
- **CLI**: Use Cobra framework, add commands in `init()` functions
- **JSON**: Use `json.NewEncoder()` for writing, `json.Unmarshal()` for reading