SHELL=/bin/bash

IMAGE=build/gobblin-bot

.PHONY: build
build:
	mkdir -p build
	gox -osarch="linux/amd64" --output="build/gobblin-bot"
	docker build -t $(IMAGE) .
	rm -rf build

.PHONY: push
push:
	docker push $(IMAGE)