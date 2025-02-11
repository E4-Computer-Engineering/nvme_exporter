#!/bin/sh
set -e

userdel -f nvme_exporter || true

systemctl daemon-reload
