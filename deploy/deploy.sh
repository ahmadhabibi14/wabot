#!/bin/bash
set -x

# Remote server information
SERVER_USER="root"
SERVER_HOST="203.194.113.211"
SSH_PORT=22

# Define SSH private key
SSH_PRIVATE_KEY="~/.ssh/habi"

# Project directory on the remote server
PROJECT_DIR="/root/wabot"

# Send project to the remote server
rsync -avz \
  --exclude=".git" \
  --exclude=".github" \
  "ssh -p ${SSH_PORT} -i ${SSH_PRIVATE_KEY}" \
  ./ $SERVER_USER@$SERVER_HOST:$PROJECT_DIR

# Run a remote script on a local machine over ssh
ssh -p $SSH_PORT -i $SSH_PRIVATE_KEY \
  $SERVER_USER@$SERVER_HOST \
  'bash' $PROJECT_DIR/deploy/build-on-server.sh