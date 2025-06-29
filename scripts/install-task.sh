#!/bin/sh
#
# This script downloads and installs the latest version of Task.
# It is designed to work on Linux, macOS, and Windows (via Git Bash or WSL).
#
# Source: https://taskfile.dev/installation

set -e

echo "🚀 Installing Task..."

# Use curl to download and execute the official installer script.
# The installer automatically detects the OS and architecture.
# -d: Download the binary to the current directory
# -b: Specify the installation directory
if command -v curl >/dev/null 2>&1; then
    sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d
else
    echo "Error: curl is required to run this script."
    exit 1
fi

# Move the downloaded binary to a common bin location
# Use sudo if necessary
if [ -w "/usr/local/bin" ]; then
    mv ./task /usr/local/bin/task
    echo "✅ Task installed successfully to /usr/local/bin/task"
else
    echo "Sudo permissions required to move Task to /usr/local/bin"
    sudo mv ./task /usr/local/bin/task
    echo "✅ Task installed successfully to /usr/local/bin/task"
fi

# Verify installation
echo ""
echo "Verifying installation..."
task --version 