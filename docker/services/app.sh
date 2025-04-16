#!/bin/bash
script_path=$(readlink -f "$0")
script_dir=$(dirname "$script_path")
proj_dir=$(dirname "$script_dir")
app_path=$(<$proj_dir/APPPATH)
if [[ $app_path == /* ]]; then
    app_path=$app_path
else
    app_path="$proj_dir/$app_path"
fi
name=$(<$app_path/NAME)
tag=$(<$app_path/VERSION)
cmd_name=${name}
container_name=${name}
image_name="$name:$tag"
network_name=zy-network
log_path=$(<$proj_dir/LOGPATH)
if [[ $log_path == /* ]]; then
    log_dir=$log_path
else
    log_dir="$proj_dir/$log_path"
fi

. "$script_dir"/helper.sh
init_network "$network_name"

usage() {
  echo "Usage: $cmd_name <command>"
  echo "$cmd_name service."
  echo "Command:"
  echo "  start, start the service"
  echo "  stop, stop the service"
  echo "  status, display service status"
  echo "  down, stop and remove the service"
  exit
}

run(){
  docker run -it \
  -v $log_dir:/var/log/app \
  --name "$container_name" \
  --restart always \
  --network "$network_name" \
  -d \
  "$image_name" > /dev/null
}

case "$1" in
start)
  if start "$container_name"; then
    run
  fi
  status "$container_name"
  ;;
restart)
  stop "$container_name"
  if start "$container_name"; then
    run
  fi
  status "$container_name"
  ;;
stop)
  stop "$container_name"
  success Done
  ;;
down)
  down "$container_name"
  success Down
  ;;
status)
  status "$container_name"
  ;;
*)
  usage
  ;;
esac

exit 0