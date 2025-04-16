FROM alpine:latest

WORKDIR /go/src/app
# 日志目录 /var/log/app
RUN mkdir -p /var/log/app
# 配置目录 /etc/app
RUN mkdir -p /etc/app

COPY app_linux_amd64 ./app
COPY config.yaml /etc/app/config.yaml

EXPOSE 8080
EXPOSE 8090

CMD ["./app", "-c", "/etc/app/config.yaml"]