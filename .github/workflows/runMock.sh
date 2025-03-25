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
# echo "Analyzing snapshot..."
# SNAPSHOT_ID=$(proxymock analyze | grep "snapshotId:" | head -n 1 | sed 's/.*snapshotId:\([^ ]*\).*/\1/')
# if [[ -z "$SNAPSHOT_ID" ]]; then
#   echo "Error: Could not extract snapshot ID from proxymock analyze output!"
#   exit 1
# fi
# echo "Using snapshot: $SNAPSHOT_ID"

# Start proxymock in the background
nohup proxymock run --service "http=18080" --service "https=18443" --dir ./proxymock > proxymock.log 2>&1 &
# Wait briefly to ensure proxymock starts
sleep 5

# Verify proxymock is running
if ! pgrep -f "proxymock run"; then
  echo "Error: Proxymock is NOT running!"
  cat proxymock.log
  exit 1
fi

echo "Proxymock started successfully."
