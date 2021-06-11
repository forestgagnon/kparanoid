.DEFAULT_GOAL := help
IMAGE_NAME_AND_TAG ?= forestgagnon/kparanoid:1

.PHONY: phony

build-cli: phony ## Build the CLI docker image
	docker build . -t $(IMAGE_NAME_AND_TAG) --file=cli.Dockerfile

# help boilerplate
BLUE := $(shell tput setaf 4)
RESET := $(shell tput sgr0)
.PHONY: help
help: ## List all targets and short descriptions of each
	@grep -E '^[^ .]+: .*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk '\
			BEGIN { FS = ": .*##" };\
			{ printf "$(BLUE)%-29s$(RESET) %s\n", $$1, $$2  }'
