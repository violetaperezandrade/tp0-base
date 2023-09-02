#!/bin/sh

ADDRESS="server"
PORT="12345"

TESTS_PASSED=0

for i in $(seq 1 10); do
    NUMBER=$((RANDOM % 100))

    RESPONSE=$(echo "$NUMBER" | nc "$ADDRESS" "$PORT")

    if [ "$RESPONSE" == "$NUMBER" ]; then
    
        TESTS_PASSED=$((TESTS_PASSED+1))
    fi
done
echo "Tests passed: $TESTS_PASSED of 10"