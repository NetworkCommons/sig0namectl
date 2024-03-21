PROJECT := sig0namectl

init:
	@./ansible-ssh-keygen.sh
	terraform init

apply:
	terraform apply -auto-approve

destroy:
	terraform destroy -auto-approve

rebuild: destroy apply
