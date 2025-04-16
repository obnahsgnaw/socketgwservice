script_dir=$(CURDIR)
proj_dir=$(shell dirname ${script_dir})
app_base_path=$(shell cat ${proj_dir}/APPPATH)
ifeq ($(findstring /,$(app_base_path)),/)
    app_dir=${app_base_path}
else
    app_dir=${proj_dir}/${app_base_path}
endif
name=$(shell cat ${app_dir}/NAME)
tag=$(shell cat ${app_dir}/VERSION)
app_path= ${app_dir}/${tag}
IMAGE_NAME=${name}:${tag}

.PHONY: help
help:
	@echo "Command:"
	@echo "  base-image: load the base docker image when offline"
	@echo "       build: build the application image"

.PHONY: base-image
base-image:
ifeq ($(shell docker images -q alpine:latest 2> /dev/null),)
	@docker load -i alpine:latest.tar
	@echo "Done"
else
	@echo "Done"
endif

.PHONY: build
build:confirm
	@cp ${app_path}/app_linux_amd64 .
	@cp ${app_path}/config.yaml .
	@if docker image inspect $(IMAGE_NAME) > /dev/null 2>&1; then docker rmi $(IMAGE_NAME); fi
	docker build -t ${IMAGE_NAME} .
	@rm ./app_linux_amd64
	@rm ./config.yaml

.PHONY: confirm
confirm:
	@read -p "Are you changed the app config.yaml for docker model [y/N]?" && [ "$$REPLY" = "y" ] && echo "Continue..." || (echo "Operation cancelled by user" && exit 1);