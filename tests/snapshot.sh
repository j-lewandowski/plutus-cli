#!/bin/bash
set -e

BINARY="./plutus_test_bin"
DB_FILE="./tests/test_db.sqlite"
OUTPUT_FILE="./tests/snapshot_before.txt"

rm -f $BINARY
rm -f $DB_FILE
rm -f $OUTPUT_FILE

echo "Building test version..."
go build -o $BINARY ./cmd/plutus

export PLUTUS_DB=$DB_FILE
echo "Using database: $PLUTUS_DB"

# --- TEST CASES ---

echo "--- 1. HELP ---" >> $OUTPUT_FILE
$BINARY help >> $OUTPUT_FILE 2>&1

echo -e "\n--- 2. STATUS ---" >> $OUTPUT_FILE
$BINARY status >> $OUTPUT_FILE 2>&1

echo -e "\n--- 3. ADD DEPOSIT (100 EUR) ---" >> $OUTPUT_FILE
$BINARY add 100 2025-12-01 >> $OUTPUT_FILE 2>&1

echo -e "\n--- 4. ADD DEPOSIT (50.50 EUR) ---" >> $OUTPUT_FILE
$BINARY add 50.50 2025-12-01 >> $OUTPUT_FILE 2>&1

echo -e "\n--- 5. STATUS ---" >> $OUTPUT_FILE
$BINARY status >> $OUTPUT_FILE 2>&1

rm -f $BINARY

echo "Done $OUTPUT_FILE."