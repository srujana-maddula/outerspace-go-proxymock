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

# Import the snapshot
echo "Analyzing snapshot..."
proxymock analyze

# Locate the correct snapshot file
FILENAME=$(find ~/.speedscale/data/snapshots -maxdepth 1 -type f -name "*.json" | head -n 1)
if [[ -z "$FILENAME" ]]; then
  echo "Error: No snapshot file found!"
  exit 1
fi

SNAPSHOT_ID=$(basename "$FILENAME" .json)
echo "Using snapshot: $SNAPSHOT_ID"

# Start proxymock in the background
nohup proxymock run --service "http=18080" --service "https=18443" --snapshot "$SNAPSHOT_ID" > proxymock.log 2>&1 &
# Wait briefly to ensure proxymock starts
sleep 5

# Verify proxymock is running
if ! pgrep -f "proxymock run"; then
  echo "Error: Proxymock is NOT running!"
  cat proxymock.log
  exit 1
fi

echo "Proxymock started successfully."
