#!/bin/bash

# Script to run proxymock replay command in a loop
# Usage: ./scripts/replay_loop.sh [sleep_seconds]

# Default sleep time in seconds
SLEEP_SECONDS=${1:-5}

# Directory containing recorded traffic
RECORDED_DIR="proxymock/recorded-2025-07-23_17-13-03.589207Z"

echo "Starting proxymock replay loop..."
echo "Sleep interval: ${SLEEP_SECONDS} seconds"
echo "Recorded directory: ${RECORDED_DIR}"
echo "Press Ctrl+C to stop"
echo ""

# Counter for tracking runs
run_count=1

while true; do
    echo "=== Run #${run_count} ==="
    
    # Alternate between two different commands
    if [ $((run_count % 2)) -eq 1 ]; then
        # Odd runs: Original command
        echo "Running: proxymock replay --in ${RECORDED_DIR} --no-out"
        proxymock replay --in "${RECORDED_DIR}" --no-out
    else
        # Even runs: Command with latency testing
        echo "Running: proxymock replay --in ${RECORDED_DIR} --no-out --for 60s --fail-if \"latency.p99 > 100\""
        proxymock replay --in "${RECORDED_DIR}" --no-out --for 60s --fail-if "latency.p99 > 100"
    fi
    
    # Check if the command was successful
    if [ $? -eq 0 ]; then
        echo "✓ Replay completed successfully"
    else
        echo "✗ Replay failed with exit code $?"
    fi
    
    echo "Sleeping for ${SLEEP_SECONDS} seconds..."
    sleep "${SLEEP_SECONDS}"
    
    # Increment counter
    ((run_count++))
    echo ""
done
