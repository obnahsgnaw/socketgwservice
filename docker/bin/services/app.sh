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
app="$app_path/$tag/app_linux_amd64"

start() {
  if [ -f /var/run/"$name".pid ]; then
    echo "The program is already running with PID $(cat /var/run/"$name".pid)"
  else
    echo "Starting the program..."
    "${app} -c ${app_path}/$tag/config.yaml" &
    sudo echo $! > /var/run/"$name".pid
    echo "The program has started with PID $!"
  fi
}

stop() {
  if [ -f /var/run/"$name".pid ]; then
    echo "Stopping the program with PID $(cat /var/run/"$name".pid)"
    kill $(cat /var/run/"$name".pid)
    sudo rm /var/run/"$name".pid
  else
    echo "The program is not running"
  fi
}

case "$1" in
  start)
    start
    ;;
  stop)
    stop
    ;;
  *)
    echo "Usage: $0 {start|stop}"
    ;;
esac

exit 0