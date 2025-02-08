# nvme_exporter

Prometheus exporter for nvme smart-log and OCP smart-log metrics

* nvme smart-log field descriptions can be found on page 209 of:
https://nvmexpress.org/wp-content/uploads/NVM-Express-Base-Specification-Revision-2.1-2024.08.05-Ratified.pdf

* nvme ocp-smart-log field descriptions can be found on page 24 of:
https://www.opencompute.org/documents/datacenter-nvme-ssd-specification-v2-5-pdf

## Build

``` bash
go build .
```

A sample Dockerfile and docker-compose.yaml are provided.

## Running

Running the exporter requires the nvme-cli package to be installed on the host.

``` bash
./nvme_exporter <flags>
```

### Flags

| Name | Description |
|----|-------------------------------------------------|
|port | Listen port number. Type: String. Default: 9998 |

## Dashboard

A sample Grafana dashboard is available:

[https://grafana.com/grafana/dashboards/14706](https://grafana.com/grafana/dashboards/14706)
