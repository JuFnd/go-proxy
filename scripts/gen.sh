#!/bin/sh

sudo chmod 777 scripts/*.sh

sudo bash scripts/gen_ca.sh
sudo bash scripts/gen_cert.sh
sudo bash scripts/install_ca.sh
