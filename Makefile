#: auto-restart golang app via makefile: https://widnyana.web.id/blog/software-engineering/automatically-restart-golang-application-on-source-code-change/
SHELL=/bin/bash -o pipefail
# PROJECT_DIR=$(shell pwd)

# Define a new Path variable into current executed subshell
export PWD := $(shell pwd)
export PATH := ${PWD}/.bin:${PATH}

BRANCH ?= $(shell git rev-parse --abbrev-ref HEAD || echo "main")
COMMIT = $(shell git log -1 --pretty=format:"%at-%h" || echo "0000000000-0000000")
COMMIT_HASH = $(shell git log -1 --pretty=format:"%h" || echo "0000000")
# replace "/" with "."
BRANCH_VERSION = $(shell echo "$(BRANCH)" | sed "s/\//./g" || echo "main")
GIT_TAG = $(shell git describe --tags --abbrev=0 2> /dev/null || echo "v0.1.0+$(BRANCH_VERSION).$(COMMIT_HASH)")
# replace backtick "`" with single quote "'"
COMMIT_MSG = $(shell git log -1 --pretty=format:"%s" | sed "s/\`/\'/g" || echo "unknown" )
BUILD_TIME = $(shell date)

# Find out which version of GO we're running
GO_VERSION = $(shell echo $(shell go version) | grep -Eo "[0-9]{1,}.[0-9]{1,}")
GO_GT_1_19 = $(shell expr $(GO_VERSION) ">=" 1.19)

# See: https://stackoverflow.com/questions/53836858/shell-function-in-gnu-makefile-results-in-unterminated-call-to-function-shel#answers-header
PROJECT_GO_PACKAGE      = $(shell cat go.mod | grep -a -h -m 1 module | sed "s/module //")
PROJECT_GO_PACKAGE_ESC  = $(shell echo $(PROJECT_GO_PACKAGE) | sed "s/\//\\\\\//g")
PROJECT_DIRS            = $(shell go list ./... | sed "s/$(PROJECT_GO_PACKAGE_ESC)\//.\//")
PROJECT_PID             = "/tmp/message-service-kata.pid"
PROJECT_APP_NAME        = $(shell basename $(PROJECT_GO_PACKAGE))
PROJECT_MAIN_ENTRYPOINT = message-service-kata
PROJECT_CMD_DIR         = ./cmd

GO_DEPENDENCIES = github.com/golangci/golangci-lint/cmd/golangci-lint \
				github.com/securego/gosec/v2/cmd/gosec \
				mvdan.cc/gofumpt \
				golang.org/x/tools/cmd/goimports \
				golang.org/x/lint/golint \
				github.com/golang/mock/mockgen

AVAILABLE_CMDS=$(shell find $(PROJECT_CMD_DIR) -type d -maxdepth 1 -mindepth 1 -type d -print0 | xargs -0 -I{} basename {})

# if the given app_name is message-service-kata, then only return PROJECT_APP_NAME as app_name
# or add PROJECT_APP_NAME as a prefix
define get_app_name
$(if $(filter $(1),$(PROJECT_MAIN_ENTRYPOINT)),$(PROJECT_APP_NAME),$(PROJECT_APP_NAME)-$(1))
endef

# This approach founds on ory Makefile
define make-go-dependency
# go install is responsible for not re-building when the code hasn't changed
bin/$2: go.sum go.mod Makefile ## Install $2 into $(PWD)/.bin/$2
	@GOBIN=$(PWD)/.bin go install $1@latest
endef
$(foreach dep, $(GO_DEPENDENCIES), $(eval $(call make-go-dependency,$(dep),$(notdir $(dep)))))

## Global variables
# Place this line at the top of your Makefile
__vars_old__ := $(.VARIABLES)

# Put this at the point where you want to see the variable values
AVAILABLE_VARS=$(foreach v, $(filter-out $(__vars_old__) __vars_old__,$(.VARIABLES)), $(if $(filter $(v),AVAILABLE_VARS),,$(v)))

GO_RUN_BUILD_FLAGS:=-ldflags="\
				-X \"$(PROJECT_GO_PACKAGE)/config.BuildTime=$(BUILD_TIME)\" \
				-X \"$(PROJECT_GO_PACKAGE)/config.CommitMsg=$(COMMIT_MSG)\" \
				-X \"$(PROJECT_GO_PACKAGE)/config.CommitHash=$(COMMIT)\" \
				-X \"$(PROJECT_GO_PACKAGE)/config.AppVersion=$(GIT_TAG)\" \
				-X \"$(PROJECT_GO_PACKAGE)/config.Branch=$(BRANCH)\" \
				-X \"$(PROJECT_GO_PACKAGE)/config.ReleaseVersion=$(BUILD_TAG)\""

go.mod: version
	@go mod download

##@ Initialize Application
.PHONY: init
init: go.mod ## Initialize the project directories after git clone
	@mkdir -p .bin
	@go install github.com/mrtazz/checkmake/cmd/checkmake@latest
	@$(if $(shell which pre-commit), \
		pre-commit install --install-hooks -t commit-msg -t pre-commit, \
		$(error No pre-commit in PATH. Ensure pre-commit is installed. See: https://pre-commit.com/#installation.))
	@echo "Successfully initialized project directory."

.PHONY: version
version: ## Validate go version
ifeq ($(GO_GT_1_19),1)
	$(info Current $(shell go version))
else
	$(error minimum supported Go version is go1.19; found $(shell go version))
endif
## -

##@ Generate
generate-mock:
	@PROJECT_DIR=${PWD} go generate ./...

##-

##@ Run rest
.PHONY: serve-rest
serve-rest:  ## Run main application and automatically restart on source code change
	go run -race $(GO_RUN_BUILD_FLAGS) $(GO_RUN_FLAGS) ./cmd/message-service-kata -service rest
## -

##@ Lint
.PHONY: verify
verify: go.mod ## Tidy & Verify Go Modules
	@go mod tidy
	@go mod verify

.PHONY: lint
lint: verify bin/goimports bin/gofumpt bin/golint bin/gosec lint-ci ## Lint this codebase
	@echo "Reformat Go import lines..."
	@goimports -v -w $(PROJECT_DIRS)
	@echo "Reformat the Go source files..."
	@gofumpt -l -w $(PROJECT_DIRS)
	@echo "Lints the Go source files..."
	@golint -set_exit_status ./...
	@echo "Check security the Go source files..."
	@gosec -stdout -severity low -no-fail $(PROJECT_DIRS)

.PHONY: lint
lint-ci: verify bin/golangci-lint ## Lint this codebase using golangci-lint
	@echo "Lints the Go source files using linter aggregator..."
	@golangci-lint run \
					--allow-parallel-runners \
					--print-resources-usage \
					--sort-results \
					--out-format colored-line-number \
					--build-tags=integration
## -

##@ Test
.PHONY: gosec
gosec: bin/gosec ## Run Gosec
	@gosec -stdout -fmt junit-xml -nosec -out gosec.xml -severity low -no-fail ./... && echo "Done!" || echo "Failed"

.PHONY: test
test: go.mod ## Run automated test
	@go test -bench -race -v -run '' -coverprofile=reports/coverage.out $(PROJECT_GO_PACKAGE)/...

.PHONY: show-coverage
show-coverage: go.mod ## Show test coverage in browser
	@go tool cover -html=reports/coverage.out
## -

##@ Clean
.PHONY: clean
clean:  ## Clean (almost) everything
	@$(foreach cmd, $(AVAILABLE_CMDS), $(shell rm -f $(call get_app_name,$(cmd)-$(BUILD_TAG))))
	@$(foreach cmd, $(AVAILABLE_CMDS), $(shell pkill $(call get_app_name,$(cmd)-$(BUILD_TAG)) 2>&1 >> /dev/null || true))
	@(MAKE) verify
	@echo "Done..."
## -

##@ Compile
.PHONY: all
all: compile compile-migration ## Compile all application
.PHONY: compile
compile: compile/message-service-kata ## Compile the application via ./cmd/message-service-kata
.PHONY: compile-migration
compile-migration: compile/migration ## Compile the migration application via ./cmd/migration

# compress: ## Compress binary file using upx
# 	@echo "Compressing binary"
# 	@upx --ultra-brute $(PROJECT_APP_NAME)
## -

.PHONY: help
.DEFAULT_GOAL := help
help:
	@printf "Usage: make \033[32m<target>\033[0m\n\n"
	@echo "[!] Available Variables:"
	@echo "------------------------"
	@printf "\033[34m%30s\033[0m  \033[36m%-s\033[0m\n" $(foreach v, $(AVAILABLE_VARS), $v $(if $($(v)),$($(v)), '-'))
	@echo ""
	@echo "[!] Available Targets:"
	@echo "----------------------"
	@awk 'BEGIN {FS = ":.*##";} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[34m%-20s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[33m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


##@ Run consumer
.PHONY: serve-consumer
serve-consumer:  ## Run main application and automatically restart on source code change
	go run -race $(GO_RUN_BUILD_FLAGS) $(GO_RUN_FLAGS) ./cmd/message-service-kata -service consumer
## -