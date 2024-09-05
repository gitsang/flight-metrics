#!/bin/bash

go build -o ./dist/bin/flight-metrics .

docker build -t flight-metrics ./

docker rm -f flight-metrics

docker run -d \
    --name flight-metrics \
    -p "2112:2112" \
    -v "${PWD}/.env:/.env" \
    flight-metrics
