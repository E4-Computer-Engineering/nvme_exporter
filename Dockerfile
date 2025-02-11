FROM ubuntu:24.04

COPY nvme_exporter /usr/bin/nvme_exporter

EXPOSE 9998
ENTRYPOINT ["/usr/bin/nvme_exporter"]
