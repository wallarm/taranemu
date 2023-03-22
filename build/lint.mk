LINT_BIN			=./.bin/golangci-lint

$(LINT_BIN): # Install linter
	@echo "Installing $@..."
	@mkdir -p $(@D)
	GOBIN=$(abspath $(@D)) go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.51.2

.PHONY: lint
lint: $(LINT_BIN) # Run linter on the codebase
	$(LINT_BIN) run ./...
