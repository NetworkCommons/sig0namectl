#!/bin/bash
export ANSIBLE_HOST_KEY_CHECKING=False
make destroy
sleep 10
sed -i.bak '/sig0namectl ssh-ed25519 /d' ~/.ssh/known_hosts
make build
sleep 100
make -C ansible sig0namectl
