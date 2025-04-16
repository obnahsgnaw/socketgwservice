export PACKAGE=github.com/obnahsgnaw/socketgwservice
export INPUT=cmd/main.go
export OUT=out
export APP=tcpgateway


.PHONY: help
help:base_help build_help

.PHONY: base_help
base_help:
	@echo "usage: make <option> <params>"
	@echo "options and effects:"
	@echo "    help   : Show help"

.PHONY: init
init:
	@go mod tidy
	@go install github.com/xiaoqidun/gitcz@v1.0.4

.PHONY: package
package:
	@go run cmd/init.go

.PHONY: asset
asset:
	@echo "Generate admin view asset file..."
	@go-bindata -o=asset/asset.go -pkg=asset html/...
	@echo "Done"

include ./build/build/makefile
include ./build/docker/makefile
include ./build/test/makefile
include ./build/version/makefile