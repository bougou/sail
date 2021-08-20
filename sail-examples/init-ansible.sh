#!/bin/bash

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

sudo yum install -y \
  ansible \
  python2-click \
  python-configparser \
  python-boto3 \
  python2-winrm \
  python2-requests_ntlm \
  python2-cryptography \
  s3cmd \
  net-tools \
  vim \
  wget \
  rsync

sudo pip install --user -r ${SCRIPT_DIR}/requirements.txt
