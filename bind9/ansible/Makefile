PROJECT := sig0namectl

play:
	ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts.yml install.yml

create-zone:
	ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts.yml -vvv create_zone.yml

so:
	ANSIBLE_HOST_KEY_CHECKING=False ansible-playbook -i hosts.yml -vvv so_playbook.yml

clean_ssh_key:
	rm ansible/ansible_ssh_key
	rm ansible/ansible_ssh_key.pub

