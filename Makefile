# Common makefile structure for buildable go projects that result in a docker container artifact
export SHELL:=/bin/bash
export SHELLOPTS:=$(if $(SHELLOPTS),$(SHELLOPTS):)pipefail:errexit

CWD=$(shell pwd)

export VERSION?=$(shell git describe --always)

# for docker image tagging and repos
export IMAGE_NAME?=outerspace-go
export REGISTRY?=ghcr.io/kenahrens
export REPO=$(REGISTRY)/$(IMAGE_NAME)

# builder image
export GOLANG_IMAGE=golang:1.23
export CGO_ENABLED?=0

# Build multi arch images if:
# 1. It isn't already defined by the user or a subproject
# 2. This is a CI master build
ifdef CI
ifndef CI_MERGE_REQUEST_ID
PLATFORMS?=linux/amd64,linux/arm64
else
PLATFORMS?=linux/amd64
endif
endif

# for local builds
PROGRAM_NAME?=$(shell basename `pwd`)
BUILD_DIR=$(CWD)/build

# buildkit builder, if set, will be used when running the build step
BUILDKIT_BUILDER?=builder

all: build

clean:
	rm -rf vendor
	rm -rf build

vendor.base:
	GOWORK=off go mod vendor

build.%: env.base vendor.base
ifdef CI
	docker context create $(BUILDKIT_BUILDER) || true
	docker buildx use builder || docker buildx create $(BUILDKIT_BUILDER) --name builder --use || true
endif
	$(eval DOCKER_TARGET=$(if $(DOCKER_BUILD_TARGET),--target=$(DOCKER_BUILD_TARGET),))
	docker buildx build $(DOCKER_TARGET) \
		$(if $(NOATTEST),,--attest=type=provenance,enabled=false) \
		-f $(CWD)/Dockerfile \
		$(if $(PLATFORMS),--platform $(PLATFORMS),) \
		$(if $(CI),--push,--output type=docker) \
		--tag $(REPO)$(IMAGE_SUFFIX):$(VERSION) \
		--build-arg VERSION=$(VERSION) \
		--build-arg BUILD_IMAGE=$(GOLANG_IMAGE) \
		--build-arg CGO_ENABLED=$(CGO_ENABLED) \
		$(CWD)

push.%: env.base
	$(eval IMAGE=$(REPO)$(IMAGE_SUFFIX))
	@echo "Pushing $(IMAGE): tag=$(VERSION) tag=latest"
# In CI, multi-arch builds get pushed during buildx, but still tag latest
ifdef CI
	@# CI: use gcloud to tag the image
	gcloud container images add-tag $(IMAGE):$(VERSION) $(IMAGE):latest --quiet
else
	@# push the version tag and add latest tag. don't push latest because CI does that
	docker push $(IMAGE):$(VERSION)
	docker tag $(IMAGE):$(VERSION) $(IMAGE):latest
endif

env.base:
ifndef IMAGE_NAME
	$(error IMAGE_NAME is not set (i.e. gcr.io/speedscale-demos/my-project))
endif

.force:

# sorcery: this is a catch-all recipe for anything inluding this file. it effectively
# allows adding functionality without getting warnings
%: .force %.base
	@# nothing
