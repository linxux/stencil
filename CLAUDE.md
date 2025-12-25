# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**Stencil** is a project scaffolding generator written in Go that creates projects from templates with variable substitution. It is language-agnostic and works with any programming language or project type.

### Core Functionality

Stencil takes a template directory and generates a new project by:
1. Recursively walking the template directory
2. Replacing variables in file contents, filenames, and directory names
3. Preserving file permissions and copying binary files without modification
4. Optionally running in interactive mode to collect variable values

## Architecture

### Module Structure

```
github.com/linxux/stencil
├── cmd/stencil/          # CLI entry point and flag handling
├── config/               # Configuration management (not internal - may be imported)
├── internal/
│   ├── generator/        # Core generation engine (walks templates, creates output)
│   ├── replacer/         # Variable substitution (handles 4 formats, binary detection)
│   └── interactive/      # User prompts and confirmation handling
└── examples/             # Template examples and usage documentation
```

### Key Components

**config/config.go**: Defines the `Config` struct and provides `LoadConfig()`, `SaveConfig()`, and `DefaultConfig()`. Note: This is in `config/` (not `internal/`) so it could be imported by other packages if needed.

**internal/generator/generator.go**: Core generation logic. The `Generator` struct holds a `*config.Config` and `*replacer.Replacer`. Key methods:
- `Generate()` - Main entry point, walks template directory
- `ExtractVariables()` - Scans template to find all variables
- `processFile()` - Handles individual file processing
- Setters (`SetVariables()`, `TemplateDir()`, etc.) - Used by CLI for state access

**internal/replacer/replacer.go**: Handles variable substitution in content and paths. Supports 4 variable formats using regex patterns:
- `{{var}}`, `<<var>>`, `__var__`, `%var%`
- Detects binary files via null byte scanning
- `ReplaceInContent()` operates on `[]byte`, `ReplaceInPath()` on strings

**internal/interactive/prompt.go**: User interaction using `bufio.Reader`. Prompts for variables, confirmation, and string input.

### Configuration Priority (Highest to Lowest)

1. Command-line flags (`-t`, `-o`, `-v`, `-i`, `--disable-*`, etc.)
2. Config file specified with `-c`
3. Auto-detected config file (`stencil.json`, `.stencil.json`, `stencil.config.json`)
4. Built-in defaults (flags defined with defaults: `-t ./template`, `-o ./output`)

**Important**: Flag defaults are set in `flag.StringVar()` calls in `init()`, not in `config.DefaultConfig()`. This means `./bin/stencil` with no args uses `./template` and `./output`.

## Format Control

### The Problem

Some variable formats can conflict with language-specific syntax. For example:
- **Go**: Uses `%s`, `%d`, etc. in `fmt.Sprintf()` which could be confused with the `%var%` format
- **Python**: Jinja2 templates use `{{var}}` which conflicts with the `{{var}}` format
- **C++**: Template syntax uses `<Type>` which could conflict with the `<<var>>` format

### The Solution

Stencil allows you to disable specific variable formats using:
1. Command-line flags: `--disable-braces`, `--disable-angle-brackets`, `--disable-underscores`, `--disable-percent`
2. Config file: `"formats": { "enableBraces": false, ... }`

All formats are **enabled by default** for backward compatibility.

### Implementation

**config/config.go**:
- `FormatOptions` struct with `EnableBraces`, `EnableAngleBrackets`, `EnableUnderscores`, `EnablePercent` bool fields
- Added to `Config` struct as `Formats FormatOptions`
- `DefaultConfig()` sets all formats to `true`

**internal/replacer/replacer.go**:
- `NewReplacer()` now accepts `formats config.FormatOptions` parameter
- `ReplaceInContent()` checks each format flag before replacing
- `ReplaceInPath()` checks each format flag before replacing
- `ExtractVariablesFromFile()` accepts `formats` parameter and only extracts enabled formats
- `ExtractVariablesFromPath()` accepts `formats` parameter and only extracts enabled formats

**internal/generator/generator.go**:
- `NewGenerator()` passes `cfg.Formats` to `NewReplacer()`
- `ExtractVariables()` passes `g.cfg.Formats` to extract functions
- `SetVariables()` passes `g.cfg.Formats` to `NewReplacer()`

**cmd/stencil/main.go**:
- Added pointer bool flags: `disableBraces`, `disableAngleBrackets`, `disableUnderscores`, `disablePercent`
- `loadConfig()` applies format flags to `cfg.Formats` (inverts boolean because flag is "disable")
- Help text updated with format flags and examples

## Development Commands

```bash
make build          # Build to ./bin/stencil with optimization flags (-ldflags="-s -w")
make dev            # Run in development mode (go run)
make test           # Run tests (currently no tests exist)
make clean          # Remove ./bin directory
make update-deps    # Update go.mod dependencies

# Manual build
go build -ldflags="-s -w" -o ./bin/stencil ./cmd/stencil/main.go
```

## Variable Substitution

### Supported Formats

All four formats are supported simultaneously in both file contents and paths:
- `{{variable}}` - Most common, clear delimiters
- `<<variable>>` - Alternative when braces conflict
- `__variable__` - Recommended for filenames/directories (e.g., `__project_name__/`)
- `%variable%` - Percent style

### Extraction Patterns

The replacer uses regex to extract variables from templates:
- `\{\{([^}]+)\}\}` for `{{var}}`
- `<<([^>]+)>>` for `<<var>>`
- `__([A-Za-z0-9_]+)__` for `__var__`
- `%([A-Za-z0-9_]+)%` for `%var%`

Variables are extracted from both file contents and paths during `ExtractVariables()`.

### Replacement Behavior

- **File contents**: Byte-level replacement using `bytes.ReplaceAll()`
- **Paths**: String-level replacement using `strings.ReplaceAll()`
- **Binary files**: Detected via null byte in first 512 bytes, copied verbatim
- **Permissions**: Preserved from template to output

## Configuration File

### Auto-Detection

When no `-c` flag is provided, Stencil checks for:
1. `stencil.json` (recommended)
2. `.stencil.json` (hidden file)
3. `stencil.config.json`

The detected file is shown: `Using config file: stencil.json`

### JSON Schema

```json
{
  "templateDir": "./template",     // Source template directory
  "outputDir": "./output",         // Generated output directory
  "variables": {                   // Variable key-value pairs
    "project_name": "myapp",
    "module_path": "github.com/example/myapp",
    "author": "Your Name",
    "description": "My app",
    "version": "1.0.0"
  },
  "interactive": false,            // Enable interactive prompts
  "dryRun": false,                 // Preview without creating files
  "skipConfirm": false             // Skip confirmation in interactive mode
}
```

### Variable Merging

Command-line variables (`-v`) are merged with config file variables. Command-line values override config values.

## Template Structure

### Example Template

```
template-go-basic/
├── README.md                  # Contains: {{project_name}}, {{description}}, {{author}}
├── go.mod                     # Contains: {{module_path}}
├── cmd/
│   └── __project_name__/      # Directory name becomes variable value
│       └── main.go           # Contains: {{project_name}}, {{module_path}}
└── internal/
    └── __project_name__/      # Package name becomes variable value
        └── app.go            # Contains: {{project_name}}, {{version}}
```

### Generation Process

1. Validate template directory exists
2. If `interactive: true`, extract variables and prompt user
3. Walk template directory recursively
4. For each path: replace variables in path
5. For directories: create with original permissions
6. For files:
   - If binary: copy verbatim
   - If text: replace variables in content, then write
7. If `dryRun: true`: show operations without executing

## Error Handling

### Template Not Found

When template directory doesn't exist, Stencil shows:

```
Error: Template directory does not exist: ./template

GETTING STARTED:

  Option 1: Use a config file (recommended)
    ...
  Option 2: Use command-line flags
    ...
  Option 3: Try the example template
    ...
```

This is implemented in `main()` after `loadConfig()` by checking `os.Stat(cfg.TemplateDir)`.

### Helpful Messages

Always provide context in error messages. Show what was attempted and suggest how to fix it. See `printGettingStarted()` for example.

## Code Conventions

### Package Organization

- `config/` - Not in `internal/`, can be imported by other packages
- `internal/` - Private implementation, cannot be imported outside this module
- Use descriptive package names: `config`, `generator`, `replacer`, `interactive`

### Error Messages

- Use `fmt.Fprintf(os.Stderr, "...")` for errors
- Include context (what operation, what file, why it failed)
- For user-facing errors, provide actionable guidance

### CLI Pattern

- Define flags in `init()` with short and long versions
- Use `flag.StringVar()` for string flags (set defaults in the call, not in variables)
- Parse flags first, then handle special cases (`--help`, `--version`)
- Load config, validate, then execute

### Interactive Mode

When `cfg.Interactive == true`:
1. Call `gen.ExtractVariables()` to find all variables in template
2. Use `prompter.PromptForValues()` to collect user input
3. Show summary with all values
4. Use `prompter.PromptForConfirmation()` unless `cfg.SkipConfirm == true`
5. Call `gen.SetVariables()` to update generator with collected values
6. Call `gen.Generate()`

## Important Implementation Details

### Binary File Detection

`replacer.IsBinaryFile()` reads first 512 bytes and checks for null bytes (`\x00`). If found, file is binary. This is reliable for most binary formats (images, executables, etc.).

### File Walking

Uses `filepath.Walk()` which walks the directory tree recursively. For each entry:
- Skip the root template directory itself (`relPath == "."`)
- Apply path replacement to get target path
- Create directories or process files based on `os.FileInfo`

### Dry Run Mode

When `cfg.DryRun == true`:
- Print `[DRY RUN]` prefix for all operations
- Show file path being created
- Show content preview (first 200 chars) for text files
- Return without actually creating files

### Variable Case Sensitivity

All variable matching is case-sensitive. `{{ProjectName}}` and `{{project_name}}` are different variables.

## Testing Notes

Currently no tests exist in the codebase. When adding tests:
- Use `*_test.go` files alongside source files
- Test both happy path and error cases
- Mock file system operations for generator tests
- Test variable extraction and replacement patterns

## Documentation

- `README.md` - User-facing documentation
- `examples/USAGE.md` - Comprehensive usage examples and tutorials
- `examples/config.json` - Example configuration file
- `examples/template-go-basic/` - Working Go project template

When updating functionality, keep these files in sync.
