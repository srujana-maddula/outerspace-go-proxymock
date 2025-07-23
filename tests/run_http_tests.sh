#!/bin/bash

# Check if recording mode is enabled via command line argument
if [ "$1" = "--recording" ] || [ "$1" = "-r" ]; then
    PORT=4143
    echo "Recording mode enabled, using port $PORT"
else
    PORT=8080
    echo "Using default port $PORT"
fi

FAIL=0
for url in $(grep '^GET' tests/test.http | sed 's/GET //' | sed "s/{{host}}/localhost:$PORT/"); do
    echo -n "Testing $url... "
    STATUS=$(curl -s -o /dev/null -w "%{http_code}" $url)
    if [ "$STATUS" -eq 200 ]; then
        echo "OK ($STATUS)"
    else
        echo "FAIL ($STATUS)"
        FAIL=1
    fi
done

if [ $FAIL -eq 1 ]; then
    echo "Http tests failed."
    exit 1
else
    echo "Http tests passed."
    exit 0
fi
