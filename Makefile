APP_NAME := ri-utilization-plotter
VETARGS?=-all
GOFMT_FILES?=$$(find . -name '*.go' | grep -v vendor)

LOGICAL_FUNCTION_NAME := RIUtilizationPlotter
S3_BUCKET := serverless.deployment.hoge

.PHONY: deps
deps:
	go get -u ./...

## Setup
.PHONY: devel-deps
devel-deps: deps
	GO111MODULE=off go get \
		golang.org/x/lint/golint \
		honnef.co/go/tools/staticcheck \
		github.com/kisielk/errcheck \
		golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow \
		github.com/securego/gosec/cmd/gosec \
		github.com/motemen/gobump/cmd/gobump \
		github.com/Songmu/make2help/cmd/make2help

.PHONY: test
test: deps
	go test -v -count=1 -cover ./...

.PHONY: cov
cov:
	go test -coverprofile=cover.out ./...
	go tool cover -html=cover.out

.PHONY: clean
clean:
	rm -rf ./dst/${APP_NAME}
	rm -rf ./dst/configs

## Lint
.PHONY: lint
lint: devel-deps
	go vet ./...
	staticcheck ./...
	errcheck ./...
	gosec -quiet ./... 
	golint -set_exit_status ./...

.PHONY: fmt
fmt:
	gofmt -s -l -w $(GOFMT_FILES)

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o dst/${APP_NAME} ./src
	cp -r ./configs ./dst

.PHONY: validate
validate:
	sam validate

.PHONY: local-invoke
local-invoke: build
	sam local invoke RIUtilizationPlotter -e testdata/event.json

.PHONY: package
package:
	sam package \
		--template-file template.yaml \
		--output-template-file packaged.yaml \
		--s3-bucket ${S3_BUCKET}

.PHONY: deploy
deploy:
	sam deploy \
		--template-file packaged.yaml \
		--stack-name ${APP_NAME} \
		--capabilities CAPABILITY_IAM

.PHONY: release
release:
	$(MAKE) clean
	$(MAKE) validate
	$(MAKE) build
	$(MAKE) package
	$(MAKE) deploy
