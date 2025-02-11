#!/bin/sh
set -e

systemctl stop nvme_exporter.service || true
systemctl disable nvme_exporter.service || true

systemctl daemon-reload
