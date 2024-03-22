PROJECT := sig0namectl

init:
	make -C terraform init

build:
	make -C terraform build

destroy:
	make -C terraform destroy

clean:
	make -C terraform clean
	make -C ansible clean_ssh_key
play:
	ANSIBLE_HOST_KEY_CHECKING=False make -C ansible play

create-zone:
	make -C ansible create-zone

so:
	make -C ansible so


all:	destroy init build play
