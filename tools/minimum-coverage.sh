#!/bin/bash
# Author: rubensayshi (see wallet-core)

MINCOVERAGE=$1
COVERAGE=$(go tool cover -func=coverage/full | grep -P 'total:\t+\(statements\)\t+([\d.]+)%' | grep -o -P '([\d.]+)')
OK=$(echo $COVERAGE '>=' $MINCOVERAGE | bc -l)

echo "$COVERAGE%"

if [ "$OK" = "1" ]; then
    echo "ok"
else
    echo "not ok"
    exit 1
fi
