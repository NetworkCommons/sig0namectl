PROJECT := sig0namectl

play:
	ansible-playbook -i hosts.yml install.yml

apply:
	terraform apply -auto-approve

create-zone:
	ansible-playbook -i hosts.yml -vvv create_zone.yml

so:
	ansible-playbook -i hosts.yml -vvv so_playbook.yml

init:
	terraform init

## recreate terraform resources
rebuild: destroy apply

destroy:
	terraform destroy -auto-approve
