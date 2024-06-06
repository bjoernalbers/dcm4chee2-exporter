# Dcm4chee 2 Prometheus Exporter

This exporter exposes metrics from (deprecated)
[dcm4chee 2](https://dcm4che.atlassian.net/wiki/spaces/ee2/overview)
for
[Prometheus](https://prometheus.io).

Currently only some metrics from the MoveScu queue are exported.

## Installation

### Download

Just download the
[latest release](https://github.com/bjoernalbers/dcm4chee2-exporter/releases/latest)
and make it executable:

    $ curl -L https://github.com/bjoernalbers/dcm4chee2-exporter/releases/latest/download/dcm4chee2-exporter-linux-amd64 -o /usr/local/bin/dcm4chee2-exporter
    $ chmod +x /usr/local/bin/dcm4chee2-exporter

### Build it yourself

Clone this repo and build the binary via `make` (requires Go to be installed).

## Usage

Start the exporter on the host running dcm4chee:

    $ ./dcm4chee2-exporter -u admin -p secret -s /opt/dcm4chee/bin/twiddle.sh
    ...

The exporter provides dcm4chee's metrics via HTTP which can be scraped by
Prometheus:

    $ curl http://localhost:9404/metrics | grep dcm4chee2
    # HELP dcm4chee2_movescu_queue_delivering_count The number of messages currently being delivered.
    # TYPE dcm4chee2_movescu_queue_delivering_count gauge
    dcm4chee2_movescu_queue_delivering_count 0
    # HELP dcm4chee2_movescu_queue_message_count The number of messages in the queue.
    # TYPE dcm4chee2_movescu_queue_message_count gauge
    dcm4chee2_movescu_queue_message_count 3
    # HELP dcm4chee2_movescu_queue_scheduled_message_count The number of scheduled messages in the queue.
    # TYPE dcm4chee2_movescu_queue_scheduled_message_count gauge
    dcm4chee2_movescu_queue_scheduled_message_count 3
    # HELP dcm4chee2_scrape_duration_seconds Duration of backend response in seconds.
    # TYPE dcm4chee2_scrape_duration_seconds gauge
    dcm4chee2_scrape_duration_seconds 0.448780636
    # HELP dcm4chee2_up Availability of Dcm4chee 2.
    # TYPE dcm4chee2_up gauge
    dcm4chee2_up 1

All available options:

    $ ./dcm4chee2-exporter --help
    Dcm4chee 2 Prometheus Exporter (version ...)

    This exporter exposes metrics from (deprecated) dcm4chee 2 for Prometheus.

    Usage:
      -a string
            Address to listen on (default ":9404")
      -p string
            Password of JBoss admin (default "admin")
      -s string
            Path to JBoss twiddle script (default "twiddle.sh")
      -u string
            Username of JBoss admin (default "admin")

    Homepage: https://github.com/bjoernalbers/dcm4chee2-exporter
