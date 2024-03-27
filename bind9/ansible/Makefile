PROJECT := sig0namectl

install:
	ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts.yml install.yml

sig0namectl:
	ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts.yml sig0namectl.yml

bind9:
	ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts.yml bind9.yml

play: install sig0namectl bind9

create-zone:
	ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts.yml -vvv create_zone.yml

so:
	ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts.yml -vvv so_playbook.yml

clean_ssh_key:
	rm ansible/ansible_ssh_key
	rm ansible/ansible_ssh_key.pub

