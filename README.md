# Project Finder

So I ended up having a bit of an issue with Project organization on my local machine. Once a folder would get too crowded I would start another Projects folder and append some sort of tag to the end of the previous directory. Example: Projects -> Projects_Q1_24.
This became hard to navigate, as I needed to remember when I was working on a particular project when it came time to revisit for a feature request or bug fix.

While learning Go, I decided it would be a cool little side project to create a CLI tool to make my life a little easier. This simple little project takes two arguments, index and find.

Index will start at a preset Directory - in this case ~/Desktop/Projects. It creates objects from all of your directories and saves it as a json file in ~/.project-finder/projects.json.

Find takes a query term and uses a fuzzy finder to compare the query to each of your directories, and returns them in an interactive terminal session to cycle and select.
They are organized by time last modified, and the time is shown to the user. 

On Select it executes VsCode cli to open or add the project to a new or current VSCode session.

## Features

- **Fuzzy Search**: Quickly find projects by typing partial names.
- **VSCode Integration**: Open projects directly in VSCode.
- **Directory Indexing**: Indexes projects in a predefined directory for fast access.

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
   go build -o project-finder
   ```

4. **Add the tool to your PATH** (optional):
   ```bash
   export PATH=$PATH:/path/to/project-finder
   ```

## Usage

### Indexing Projects

The first thing you need to do is set the config.
For now the CLI only accepts three directories:
```go
Items: []string{"Desktop", "Documents", "Downloads"}
```
In the future I will expand this to dynamically search through directories to add. I also want to enable the ability to add multiple directories,
but in theory you should really only have to add one parent directory.

To run the config selection:
```bash
findit config
```


To index all projects in the default directory (`~/Desktop/Projects`), run:

```bash
project-finder index
```

This command will create a `projects.json` file in the `.project-finder` directory in your home folder.

### Finding and Opening Projects

To find and open a project, use the `find` command followed by your search query:

```bash
project-finder find <name-of-directory>
```
