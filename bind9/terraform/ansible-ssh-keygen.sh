#!/bin/bash


export ANSIBLE_SSH_KEY="../ansible/ansible_ssh_key"
[ ! -f ${ANSIBLE_SSH_KEY} ] && ssh-keygen -f ${ANSIBLE_SSH_KEY} -t ed25519 -q -N "" -C "ansible-ssh-key" || echo "${ANSIBLE_SSH_KEY} already exists ... updating cloud_init.cfg"
sed -i 's/\(.*\)ssh-ed25519\(.*\)/echo '\''\1'\''$(cat ${ANSIBLE_SSH_KEY}.pub)/e' cloud_init.cfg
# cat cloud_init.cfg
