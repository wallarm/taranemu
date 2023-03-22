.DEFAULT_GOAL := help

.PHONY: help
help: # Show this help
	@egrep -h '\s#\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?# "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

include ./build/init.mk
include ./build/build.mk
include ./build/lint.mk
