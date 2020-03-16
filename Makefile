## Stack Name
STACK_NAME := ri-utilization-plotter

## Install library for production
deps:
	go get -u ./...
.PHONY: deps

## Install library for development
devel-deps: deps
	GO111MODULE=off go get \
		golang.org/x/lint/golint \
		honnef.co/go/tools/staticcheck \
		github.com/kisielk/errcheck \
		golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow \
		github.com/securego/gosec/cmd/gosec \
		github.com/motemen/gobump/cmd/gobump \
		github.com/Songmu/make2help/cmd/make2help
.PHONY: devel-deps

## Execute unit test
test: deps
	go test -v -count=1 -cover ./...
.PHONY: test

## Output test coverage
cov:
	go test -count=1 -coverprofile=cover.out ./...
	go tool cover -html=cover.out
.PHONY: cov

## Clean up artifact
clean:
	rm -rf ./artifact/*
.PHONY: clean

## Lint
lint: devel-deps
	go vet ./...
	staticcheck ./...
	errcheck ./...
	gosec -quiet ./...
	golint -set_exit_status ./...
.PHONY: lint

## Build
build: build-ri-utilization-plotter
.PHONY: build

## Build ri-utilization-plotter
build-ri-utilization-plotter:
	GOOS=linux GOARCH=amd64 go build -o artifact/ri-utilization-plotter ./handlers/ri-utilization-plotter
	cp -r ./configs ./artifact
.PHONY: build-ri-utilization-plotter

## SAM Validate
validate:
	sam validate
.PHONY: validate

## Invoke Lambda Function in local by SAM
local-invoke: build
	sam local invoke RIUtilizationPlotter -e testdata/event.json
.PHONY: local-invoke

## Deploy by SAM
deploy: build
	sam deploy
.PHONY: deploy

## Delete CloudFormation Stack
delete:
	aws cloudformation delete-stack --stack-name $(STACK_NAME)
.PHONY: delete

help:
	@make2help $(MAKEFILE_LIST)
.PHONY: help
