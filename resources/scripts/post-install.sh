#!/bin/sh
set -e

PID1=$(ps --no-headers -o comm 1)

if [ "$PID1" != systemd ]; then
	echo "Only systemd is supported but detected pid 1: $PID1"
	exit 1
fi

echo "Detected systemd as init system, proceeding"

useradd -r nvme_exporter -s /bin/false || true

systemctl daemon-reload

systemctl enable nvme_exporter.service
systemctl restart nvme_exporter.service
