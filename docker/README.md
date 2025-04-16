1. app目录中存放打包的可执行文件以及名称和版本，需要按需调整 `NAME` 和 `VERSION` 以及 `config.yaml` 可以通过 `APPPATH` 来调整路径
2. build目录中存放的是进行docker构建的脚本
3. log用来存放日志， 可以通过 `LOGPATH` 来调整路径
4. 使用
   1. docker 
      - `make init` 初始化 `docker` 脚本
      - `make auto` 将 `docker`脚本 加入系统服务
      - `make image` 构建 docker 镜像
   2. bin
      - `make init` 初始化 `docker` 脚本
      - `make auto` 将 `docker`脚本 加入系统服务
   3. 使用
      - 假设 app/NAME 中存放的是 zy-user 
      - `$ zy-user start`
      - `$ zy-user stop`
      - `$ zy-user status`
      - 如果添加了系统服务：
      - `$ systemctl start zy-user.service`
      - `$ systemctl stop zy-user.service`
5. 更新
   - 替换app目录下的可执行文件和VERSION即可 （如果docker运行需要重新`make build`）