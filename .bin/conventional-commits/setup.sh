#!/bin/bash

set -e

echo "Setting up commit-msg hook..."

# Should be run from the root of the repository
if [ ! -d ".git" ]; then
  echo "Please run this script from the root of the repository."
  exit 1
fi

# Check if git is installed
if ! command -v git &> /dev/null; then
  echo "Git is not installed. Please install Git."
  exit 1
fi

# Check if the commit-msg hook already exists
if [ -f ".git/hooks/commit-msg" ]; then
  echo "Commit-msg hook already exists. Skipping setup."
  exit 0
fi

# Define the commit-msg hook file path
HOOK_DIR=".git/hooks"
HOOK_FILE="$HOOK_DIR/commit-msg"

# Create the hooks directory if it doesn't exist
if [ ! -d "$HOOK_DIR" ]; then
  echo "Creating hooks directory..."
  mkdir -p "$HOOK_DIR"
fi

# Create or overwrite the commit-msg hook
cat << 'EOF' > "$HOOK_FILE"
#!/bin/bash

# 1. A concise summary of up to 72 characters.
# 2. An optional body that can have any length of text after a blank line.
commit_regex="^((Merge[ a-z-]* branch.*)|(Revert.*)|((build|chore|ci|docs|feat|fix|perf|refactor|revert|style|test)(\(.*\))?!?: .{1,72}(\n\n.*)?))"

# Get the commit message
commit_message=$(cat "$1")

echo "Commit message:"
echo "$commit_message"
echo ""

# Check if the commit message matches the regex
if ! echo "$commit_message" | grep -qE "$commit_regex"; then
    echo "ERROR: Commit message does not conform to the required format."
    echo "Please use the format: <type>[optional scope]: <description>"
    echo ""
    echo "Allowed types:"
    echo "  feat     : A new feature"
    echo "  fix      : A bug fix"
    echo "  chore    : Changes to the build process or auxiliary tools"
    echo "  docs     : Documentation only changes"
    echo "  style    : Changes that do not affect the meaning of the code (white-space, formatting, etc.)"
    echo "  refactor : A code change that neither fixes a bug nor adds a feature"
    echo "  perf     : A code change that improves performance"
    echo "  test     : Adding missing tests or correcting existing tests"
    echo "  build    : Changes that affect the build system or external dependencies"
    echo "  ci       : Changes to our CI configuration files and scripts"
    echo "  revert   : Reverts a previous commit"
    exit 1
fi
EOF

# Make the hook executable
chmod +x "$HOOK_FILE"

echo "$HOOK_FILE is created."
echo "Commit-msg hook has been set up successfully!"