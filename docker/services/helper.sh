#!/bin/bash

start(){
  if ! exist "$1"; then
    return 0
  else
    if ! running "$1";then
      docker start "$1" > /dev/null
    fi
    return 1
  fi
}

stop(){
  if exist "$1" && running "$1"; then
      docker stop "$1" > /dev/null
  fi
}

down(){
  if exist "$1" ; then
    docker stop "$1" > /dev/null
    docker rm "$1" > /dev/null
  fi
}

status(){
  if exist "$1"; then
    if running "$1"; then
      docker ps -a --format '{{.Names}} --- {{.Status}}' | grep "$1"
    else
      echo "stopped"
    fi
  else
    echo "not running"
  fi
}

exist(){
  if docker ps -a --format '{{.Names}}' | grep -q "^$1\$"; then
    return 0
  else
    return 1
  fi
}

running() {
  container_state=$(docker inspect -f '{{.State.Status}}' $1 2>/dev/null)
  if [ "$container_state" = "running" ]; then
    return 0
  else
    return 1
  fi
}

network_exist(){
  nk=$1
  if docker network ls --format '{{.Name}}' | grep -q "$nk"; then
      return 0
  else
      return 1
  fi
}
init_network(){
  nk1=$1
  if ! network_exist "$nk1"; then
    docker network create "$nk1" >/dev/null
  fi
}

isMac(){
  uNames=$(uname -s)
  osName=${uNames:0:4}

  if [ "$osName" == "Darw" ]; then
    return 0
  else
    return 1
  fi
}

replace(){
  search=$1
  place=$2
  file=$3
  if isMac; then
    sed -i '' "s#$search#$place#" "$file"
  else
    sed -i "s#$search#$place#" "$file"
  fi
}

success() {
  echo -e "\\033[32m $1 \\033[0m"
}