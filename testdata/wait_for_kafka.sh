#!/bin/bash

set -euo pipefail

if [ "$#" -ne 3 ]; then
    echo "Illegal number of parameters"
    echo "Usage: ./wait_for_kafka.sh <kafka_service_name> <broker_url> <broker_port> (from inside docker)"
    echo "Override settings with env-vars: MAX_TIME=10 INTERVAL=1 ./wait_for_kafka.sh x y z"
    exit 1
fi

SERVICE=$1
BROKER_URL=$2
BROKER_PORT=$3

printf 'Checking if service is up %s' "$SERVICE"

MAX_TIME=${MAX_TIME:-30}
INTERVAL=${INTERVAL:-3}
i=0


until docker exec "$SERVICE" kafka-broker-api-versions.sh --bootstrap-server "$BROKER_URL:$BROKER_PORT" > /dev/null 2>&1; do
    printf '.'
    sleep "$INTERVAL"
    ((i=i+INTERVAL))

    if (( i > MAX_TIME )); then
        printf '\nwaiting for %s timed out \n' "$SERVICE"
        exit 1
    fi
done

printf '\n%s should be up\n' "$SERVICE"