
GOCMD=go
GOTEST=$(GOCMD) test
GOVET=$(GOCMD) vet
BINARY_NAME=Exercise
MOD_BOMFILE_NAME=Exercise-bom.json	# Define the SBOM file
BENCH_TIME=1s	# Set the bench time for go bench tests, default 1s
LINTCMD=golangci-lint
LINT_CONFIG_FILE=.git/configs/golangci.yaml	# set the lint file correctly
EMPTY :=
DOCCMD=godoc-static
DOC_SITE_NAME=Exercise
DOC_SITE_FOOTER=Exercise
DOC_FOLDER=server/docs
SBOMCMD=cyclonedx-gomod


export CONTAINER_CLI=nerdctl

# import env file
env-file ?= .env
# Conditional include .env
-include env-file

ifdef OS
    RM = del /Q /F
    MKDIRP = mkdir $(subst /,\,$1)
    AWK = gawk
else
    RM = rm -rf
    MKDIRP = mkdir -p $1
    AWK = awk
endif


.PHONY: all test build

all: help

## Build:
generate: ## Generate project code
	$(GOCMD) generate ./...


build: ## Build your project and put the output binary in bin/
	-$(call MKDIRP, bin)
	$(GOCMD) build -o bin/$(BINARY_NAME) .


clean: ## Remove build generated files
	$(RM) bin
	$(GOCMD) clean

## Docs:
doc: ## Build Go docs for the project
	-$(RM) $(DOC_FOLDER)
	-$(call MKDIRP, $(DOC_FOLDER))
	$(DOCCMD) --site-name=$(DOC_SITE_NAME) --site-footer=$(DOC_SITE_FOOTER) --destination=$(DOC_FOLDER)

## Test:
test: unit-test bench-test test-coverage ## Run the unit tests, bench tests and test coverage

unit-test: ## Run the unit tests of the project
	$(GOTEST) -v ./...

bench-test: ## Run the bench tests of the project
	$(GOTEST) -bench=. ./... -v -benchmem -benchtime=$(BENCH_TIME)

test-coverage: ## Run the tests of the project and export the coverage
	$(GOTEST) -cover -covermode=count -coverprofile=profile.cov ./...
	$(GOCMD) tool cover -func profile.cov

## Lint:
lint: lint-go ## Run all lints, currently only lint-go

lint-go: ## Use golintci-lint on your project
	$(LINTCMD) run ./... -v --config=$(LINT_CONFIG_FILE)

## Run:
run: ## Run the component with go run main.go
	go run main.go

build-run: ## Run the go native binary
	-$(call MKDIRP, bin)
	$(GOCMD) build -o bin/$(BINARY_NAME) .
	./bin/$(BINARY_NAME)

run-container-it: ## Run the container interactively
	$(CONTAINER_CLI) run -it --rm --name="$(BINARY_NAME)" $(BINARY_NAME)

run-container: ## Run the go container in detached mode
	$(CONTAINER_CLI) run -d --env-file=$(env-file) --name="$(BINARY_NAME)" $(BINARY_NAME)

## Container:
container-build: ## Use the dockerfile to build the container
	$(CONTAINER_CLI) build --rm --tag $(BINARY_NAME) .

## SBOM:
sbom-gomod: ## Use cyclonedx-gomod to generate sbom for go modules
	$(SBOMCMD) mod -json -output $(MOD_BOMFILE_NAME)

## Help:
help: ## Show this help.
	$(info )
	$(info Usage:)
	$(info $(EMPTY)  make <target>)
	$(info )
	$(info Targets)
	@$(AWK) 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
