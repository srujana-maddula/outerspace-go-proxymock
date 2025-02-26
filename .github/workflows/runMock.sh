#!/bin/bash
set -e  # Exit on error
set -o pipefail  # Ensure failures in pipes are caught

# Ensure the .speedscale directory exists
mkdir -p .speedscale

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
proxymock init --api-key "$PROXYMOCK_API_KEY"

# Verify installation
proxymock version || { echo "Proxymock installation failed"; exit 1; }

# Find all JSON files inside .speedscale and its subdirectories, then concatenate them into .speedscale/raw.jsonl
find .speedscale -type f -name "*.json" -exec cat {} + | jq -c '.' > .speedscale/raw.jsonl
echo "Combined JSON files into .speedscale/raw.jsonl"

# Import and run proxymock
echo "Importing snapshot..."
proxymock import --file .speedscale/raw.jsonl

# Run proxymock with the extracted snapshot ID
FILENAME=$(ls ~/.speedscale/data/snapshots/*.json)
SNAPSHOT_ID=$(basename "$FILENAME" .json)
echo "Running proxymock with snapshot ID $SNAPSHOT_ID..."
nohup ~/.speedscale/proxymock run --snapshot-id "$SNAPSHOT_ID" --service http=18080 > proxymock.log 2>&1 &

echo "Proxymock started successfully."
