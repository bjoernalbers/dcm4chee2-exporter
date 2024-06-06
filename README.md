# Dcm4chee 2 Prometheus Exporter

This exporter exposes metrics from (deprecated)
[dcm4chee 2](https://dcm4che.atlassian.net/wiki/spaces/ee2/overview)
for
[Prometheus](https://prometheus.io).

## Installation

### Download

Just download the
[latest release](https://github.com/bjoernalbers/dcm4chee2-exporter/releases/latest)
and make it executable:

    $ curl -L https://github.com/bjoernalbers/dcm4chee2-exporter/releases/latest/download/dcm4chee2-exporter-linux-amd64 -o /usr/local/bin/dcm4chee2-exporter
    $ chmod +x /usr/local/bin/dcm4chee2-exporter

### Build it yourself

Clone this repo and build the binary via `make` (requires Go to be installed).
