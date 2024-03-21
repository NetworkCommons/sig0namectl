#!/bin/bash
export ANSIBLE_HOST_KEY_CHECKING=False
make destroy
sleep 10
make build
sleep 100
make play
#make create-demouser
