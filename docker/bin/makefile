script_dir=$(CURDIR)
project_dir=$(shell dirname ${script_dir})
app_path=$(shell cat ${project_dir}/APPPATH)
ifeq ($(findstring /,$(app_path)),/)
    app_dir=${app_path}
else
    app_dir=${project_dir}/${app_path}
endif
name=$(shell cat ${app_dir}/NAME)
cmd_name=${name}
log_path=$(shell cat ${project_dir}/LOGPATH)
ifeq ($(findstring /,$(log_path)),/)
    log_dir=${log_path}
else
    log_dir=${project_dir}/${log_path}
endif
script_name=app.sh
service_name=${cmd_name}.service
bin_file=/usr/local/bin/${cmd_name}
cmd_file=${script_dir}/services/${script_name}
service_file=${script_dir}/services/${service_name}

PLATFORM_WIN=windows
PLATFORM_MAC=darwin
PLATFORM_LINUX=linux
ifeq ($(OS),Windows_NT)
 PLATFORM=${PLATFORM_WIN}
else
 ifeq ($(shell uname),Darwin)
  PLATFORM=${PLATFORM_MAC}
 else
  PLATFORM=${PLATFORM_LINUX}
 endif
endif

.PHONY: help
help:
	@echo "Command:"
	@echo "  init: init the cmd"
	@echo "  auto: add service to startup"

.PHONY: init
init:confirm confirm-c
	@[ -f ${bin_file} ] && sudo rm -f ${bin_file}
	@sudo chmod +x ${cmd_file} && sudo ln -s ${cmd_file} ${bin_file} && sudo chmod +x ${bin_file}
	@[ -d /var/log/app ] || sudo ln -s ${log_dir} /var/log/app
	@echo "Done"

.PHONY: auto
auto:confirm
	@cp ${script_dir}/services/bin/app.service ${service_file}
ifeq (${PLATFORM},${PLATFORM_MAC})
	@sed -i '' 's/_app_/${cmd_name}/g' ${bin_service_file}
else
	@sed -i 's/_app_/${cmd_name}/g' ${bin_service_file}
endif
	@sudo mv -f ${service_file} /usr/lib/systemd/system/ >/dev/null 2>&1 && sudo systemctl enable ${service_name}
	@echo "Done"

.PHONY: confirm
confirm:
	@if [ "$(cmd_name)" = "zy-framework" ]; then \
  		read -p "Are you sure you want to continue with the cmd_name=framework [y/N]?" && [ "$$REPLY" = "y" ] && echo "Continue..." || (echo "Operation cancelled by user" && exit 1); \
	fi

.PHONY: confirm-c
confirm-c:
	@read -p "Are you changed the app config.yaml for bin model [y/N]?" && [ "$$REPLY" = "y" ] && echo "Continue..." || (echo "Operation cancelled by user" && exit 1);