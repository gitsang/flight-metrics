#!/bin/bash

docker rm -f flight-metrics

docker run -d \
    --name flight-metrics \
    --restart always \
    -p "2112:2112" \
    -v "${PWD}/.env:/.env" \
    flight-metrics
