#!/bin/bash
set -e

BINARY="heyad"
OUTPUT="${1:-genesis.json}"

print_step() { echo -e "\n\e[1;34m>>> $1\e[0m"; }

print_step "Eksportowanie genesis z dzialajacego lancucha..."
$BINARY export --home ~/.heya > "$OUTPUT"

print_step "Genesis wyeksportowany do: $OUTPUT"
echo "Skopiuj go do katalogu zrodlowego: cp $OUTPUT /root/heya/genesis.json"
