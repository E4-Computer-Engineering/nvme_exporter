FROM ubuntu:24.04

RUN apt-get update
    && apt-get install -y nvme-cli
    && rm -rf /var/lib/apt/lists/*

COPY nvme_exporter /usr/bin/nvme_exporter

EXPOSE 9998
ENTRYPOINT ["/usr/bin/nvme_exporter"]
