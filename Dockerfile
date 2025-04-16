# step 1
FROM golang:1.19 AS build

RUN apt-get update && apt-get install -y make git

RUN mkdir -p /go/src/app
COPY . /go/src/app
WORKDIR /go/src/app

ENV GOPROXY=https://goproxy.io,direct

RUN make docker-install

# setp 2
FROM alpine:latest

WORKDIR /go/src/app
# 日志目录 /var/log/app
RUN mkdir -p /var/log/app
# 配置目录 /etc/app
RUN mkdir -p /etc/app

COPY --from=build /go/src/app/app ./
COPY config.example.yaml /etc/app/config.yaml

EXPOSE 8080
EXPOSE 8090

CMD ["./app", "-c", "/etc/app/config.yaml"]