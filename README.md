# nvme_exporter

Prometheus exporter for nvme smart-log and OCP smart-log metrics

Specification versions of reference:

* nvme smart-log field descriptions can be found on page 209 of:
https://nvmexpress.org/wp-content/uploads/NVM-Express-Base-Specification-Revision-2.1-2024.08.05-Ratified.pdf

* nvme ocp-smart-log field descriptions can be found on page 24 of:
https://www.opencompute.org/documents/datacenter-nvme-ssd-specification-v2-5-pdf

Supported [NVMe CLI](https://github.com/linux-nvme/nvme-cli) versions:

| Version | Supported |
|----|----|
|2.9 | OK |
|2.10 | TBD |
|2.11 | TBD |

## Content

* Docker: A sample Dockerfile and docker-compose.yaml are provided.
* Kubernetes: In [resources](resources/k8s/).
* Grafana: In [resources](resources/grafana/) for dashboards.
  * [smart-log dashboard](https://grafana.com/grafana/dashboards/14706)
* Prometheus: In [resources](resources/prom/) for recording and alert rules.

## Running

Running the exporter requires the nvme-cli package to be installed on the host.

``` bash
./nvme_exporter -h
```

### Flags

| Name | Description | Default |
|----|----|----|
|port | Listen port number. Type: String. | `9998` |
|ocp | Enable OCP smart log metrics. Type: Bool. | `false` |
|endpoint | The endpoint to query for metrics. Type: String. | `/metrics` |

## Build

``` bash
go build .
```
