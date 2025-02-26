#!/bin/bash
set -e  # Exit on error
set -o pipefail  # Ensure failures in pipes are caught

# Ensure the .speedscale directory exists
mkdir -p .speedscale

# Find all JSON files inside .speedscale and its subdirectories, then concatenate them into .speedscale/raw.jsonl
find .speedscale -type f -name "*.json" -exec cat {} + | jq -c '.' > .speedscale/raw.jsonl
echo "Combined JSON files into .speedscale/raw.jsonl"

# Install proxymock
echo "Installing proxymock..."
sh -c "$(curl -Lfs https://downloads.speedscale.com/proxymock/install-proxymock)"
echo "Proxymock installed successfully."

# Add proxymock installation directory to PATH
export PATH="$HOME/.speedscale:$PATH"
echo "Updated PATH to include proxymock: $PATH"

# Initialize proxymock with API key
if [[ -z "$PROXYMOCK_API_KEY" ]]; then
  echo "Error: PROXYMOCK_API_KEY is not set."
  exit 1
fi

echo "Initializing proxymock..."
~/.speedscale/proxymock init --api-key "$PROXYMOCK_API_KEY"
# Verify installation
proxymock version || { echo "Proxymock installation failed"; exit 1; }

# Import and run proxymock
echo "Importing snapshot..."
~/.speedscale/proxymock import --file .speedscale/raw.jsonl

# Find the imported snapshot file
FILENAME=$(ls ~/.speedscale/data/snapshots/*.json)
echo "Snapshot filename: ${FILENAME}"

# Extract the snapshot ID
SNAPSHOT_ID=$(basename "$FILENAME" .json)
echo "Using snapshot: $SNAPSHOT_ID"

# Run proxymock with the extracted snapshot ID
echo "Running proxymock with snapshot ID $SNAPSHOT_ID..."
~/.speedscale/proxymock run --snapshot-id "$SNAPSHOT_ID"

echo "Proxymock run completed successfully."
