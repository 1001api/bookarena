#!/bin/bash
# Watch for changes in .templ files, regenerate templates, and reload browser.
# Default proxy: http://localhost:8181
# Default live URL: http://localhost:7331

# Exit if any command fails
set -e

echo "🚀 Starting Templ live generation..."
echo "Proxy: http://localhost:8181"
echo "Opening browser at http://localhost:7331"

# Run templ in watch mode with browser reload
templ generate --watch --proxy="http://localhost:8181" --open-browser=true