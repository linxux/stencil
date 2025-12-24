# Stencil Usage Examples

This document provides practical examples of using Stencil for various scenarios.

## Table of Contents

1. [Quick Start](#quick-start)
2. [Basic Usage](#basic-usage)
3. [Interactive Mode](#interactive-mode)
4. [Configuration Files](#configuration-files)
5. [Creating Your Own Templates](#creating-your-own-templates)
6. [Advanced Examples](#advanced-examples)

## Quick Start

```bash
# Build the tool
make build

# Run with example template (interactive mode)
./bin/stencil -t ./examples/template-go-basic -o ./my-project -i

# Run with variables directly
./bin/stencil -t ./examples/template-go-basic -o ./my-project \
  -v "project_name=myapp,module_path=github.com/user/myapp,author=Your Name"
```

## Basic Usage

### Generate a Go Project

```bash
./bin/stencil \
  -t ./examples/template-go-basic \
  -o ./output/myapp \
  -v "project_name=myapp" \
  -v "module_path=github.com/example/myapp" \
  -v "author=John Doe" \
  -v "description=My awesome app" \
  -v "version=1.0.0"
```

### Dry Run (Preview)

```bash
./bin/stencil \
  -t ./examples/template-go-basic \
  -o ./output/myapp \
  --dry-run \
  -v "project_name=test"
```

## Interactive Mode

Interactive mode automatically detects all variables in your template and prompts you for values:

```bash
./bin/stencil -t ./examples/template-go-basic -o ./my-project -i
```

Output:
```
=== Stencil - Interactive Mode ===
Scanning template for variables...
Found 5 variables in template.

[1/5] project_name: myapp
[2/5] module_path: github.com/example/myapp
[3/5] author: John Doe
[4/5] description: An awesome application
[5/5] version (default: 1.0.0):

=== Summary ===
Template: ./examples/template-go-basic
Output: ./my-project

Variables:
  project_name = myapp
  module_path = github.com/example/myapp
  author = John Doe
  description = An awesome application
  version = 1.0.0

Proceed with generation? [y/N]: y

Generating project...
✓ Project generated successfully!
```

## Configuration Files

### Create a Configuration File

`my-config.json`:
```json
{
  "templateDir": "./examples/template-go-basic",
  "outputDir": "./output/myapp",
  "variables": {
    "project_name": "myapp",
    "module_path": "github.com/example/myapp",
    "author": "Your Name",
    "description": "My application",
    "version": "1.0.0"
  },
  "interactive": false,
  "dryRun": false
}
```

### Use Configuration File

```bash
./bin/stencil -c my-config.json
```

### Override Config with Command-Line

```bash
# Use config but override output directory
./bin/stencil -c my-config.json -o ./different-output

# Use config but enable dry-run
./bin/stencil -c my-config.json --dry-run
```

## Creating Your Own Templates

### Template Structure

```
my-template/
├── README.md
├── package.json
├── src/
│   └── __project_name__/
│       ├── index.js
│       └── components/
│           └── App.jsx
└── config/
    └── __project_name__.config.js
```

### Supported Variable Formats

- `{{variable}}` - Double curly braces
- `<<variable>>` - Double angle brackets
- `__variable__` - Double underscores (best for filenames)
- `%variable%` - Percent signs

### Example Template Files

#### README.md
```markdown
# {{project_name}}

{{description}}

## Installation
```bash
npm install {{project_name}}
```

## Author
{{author}}
```

#### package.json
```json
{
  "name": "{{project_name}}",
  "version": "{{version}}",
  "description": "{{description}}",
  "author": "{{author}}"
}
```

#### src/__project_name__/index.js
```javascript
const PROJECT_NAME = '{{project_name}}';
const VERSION = '{{version}}';

console.log(`Welcome to ${PROJECT_NAME} v${VERSION}`);
```

## Advanced Examples

### Web Application Template

Create a React app template:

```bash
./bin/stencil \
  -t ./templates/react-app \
  -o ./my-react-app \
  -v "app_name=MyReactApp,author=John,port=3000"
```

### Microservice Template

Generate a Go microservice:

```bash
./bin/stencil \
  -t ./templates/go-microservice \
  -o ./services/user-service \
  -v "service_name=user-service,service_port=8080,db_name=users_db"
```

### Multi-Project Setup

Generate multiple projects from the same template:

```bash
# Generate service A
./bin/stencil -t ./microservice-template -o ./services/auth \
  -v "service_name=auth,port=8001"

# Generate service B
./bin/stencil -t ./microservice-template -o ./services/users \
  -v "service_name=users,port=8002"
```

### Skip Confirmation in Scripts

```bash
# In a script, use -y to skip confirmation
./bin/stencil -t ./template -o ./output -i -y
```

## Tips and Best Practices

1. **Use `__variable__` for filenames**: This is the most readable format for file and directory names

2. **Organize templates logically**: Group related files in directories

3. **Document your templates**: Include a README in your template directory explaining the variables

4. **Use dry-run first**: Always test with `--dry-run` to preview changes

5. **Version your templates**: Keep templates in version control for reproducibility

6. **Create config files for teams**: Share configuration files for consistent project generation

7. **Binary files**: Binary files (images, PDFs, etc.) are automatically detected and copied without modification

## Troubleshooting

### Variable Not Replaced

Make sure you're using one of the supported formats:
- ✅ `{{project_name}}`
- ✅ `__project_name__`
- ✅ `<<project_name>>`
- ✅ `%project_name%`
- ❌ `{project_name}` (single braces - not supported)

### File Permission Issues

Stencil preserves file permissions from the template. If you need different permissions, adjust them after generation or fix the template file permissions.

### Large Templates

For large templates, use interactive mode with `-y` flag to skip confirmation after reviewing the summary.

## Getting Help

```bash
# Show help
./bin/stencil --help

# Show version
./bin/stencil --version
```
