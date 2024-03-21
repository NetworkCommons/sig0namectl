PROJECT := sig0namectl

init:
	make -C terraform init

build:
	make -C terraform apply

destroy:
	make -C terraform destroy

play:
	make -C ansible play

create-zone:
	make -C ansible create-zone

so:
	make -C ansible so

deploy: destroy build play

all:	destroy init build play
