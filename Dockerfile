FROM centos:7

COPY dist/bin/flight-metrics /flight-metrics

ENTRYPOINT ["/flight-metrics"]
