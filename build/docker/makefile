.PHONY: docker-install
docker-install:
	@CGO_ENABLED=0 go build -o ./app -ldflags "${LDFLAGS}" ${INPUT_NAME}

IMAGE_NAME=zy-${APP_NAME}:${VERSION}

.PHONY: docker
docker:
	@if docker image inspect ${IMAGE_NAME} > /dev/null 2>&1; then docker rmi ${IMAGE_NAME}; fi
	@docker build -t ${IMAGE_NAME} .

.PHONY: image
image:
ifeq (${PLATFORM},${PLATFORM_MAC})
	@if ! docker buildx ls | grep 'myBuilder' ; then docker buildx create --config=${CURDIR}/build/docker/buildkitd.toml --name myBuilder && docker buildx inspect --bootstrap; fi
	@docker buildx build --platform linux/amd64 -t 123.57.246.85:5000/library/${IMAGE_NAME} . --push
else
ifeq (${PLATFORM},${PLATFORM_WIN})
	@echo "window build not support now..."
else
	@if docker image inspect ${IMAGE_NAME} > /dev/null 2>&1; then docker rmi ${IMAGE_NAME}; fi
	@docker build -t 123.57.246.85:5000/library/${IMAGE_NAME} . --push
endif
endif
