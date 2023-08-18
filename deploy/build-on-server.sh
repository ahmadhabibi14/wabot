#!/bin/bash
set -x

SERVICE_FILE="whatsappbot.service"
SERVICE_DIR="/lib/systemd/system"

# Go to project directory
cd /root/wabot

# Install modules and build Golang
go get
go build main.go
# Check if there is a systemd service for golang program
if [ -f "${SERVICE_DIR}/${SERVICE_FILE}" ]; then
  echo "The service file '${SERVICE_FILE}' exists in ${SERVICE_DIR}"
  systemctl restart $SERVICE_FILE
else
  echo -e \
    "\
      The file '${SERVICE_FILE}' does not exist in ${SERVICE_DIR} \n\
      SystemD service will be created now\
    "
  cp ${SERVICE_FILE} ${SERVICE_DIR}/
  systemctl daemon-reload
  systemctl start $SERVICE_FILE
fi