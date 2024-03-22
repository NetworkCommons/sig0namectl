PROJECT := sig0namectl

init:
	@./ansible-ssh-keygen.sh
	terraform init

build:
	terraform apply -auto-approve

destroy:
	terraform destroy -auto-approve

clean:
	terraform destroy -auto-approve
	rm -fr .terraform
	rm .terraform.lock.hcl

rebuild: destroy init apply
