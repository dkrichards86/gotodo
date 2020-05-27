PKGS     = $(or $(PKG),$(shell env GO111MODULE=on $(GO) list ./...))
TESTPKGS = $(shell env GO111MODULE=on $(GO) list -f \
			'{{ if or .TestGoFiles .XTestGoFiles }}{{ .ImportPath }}{{ end }}' \
			$(PKGS))
BIN		 = $(CURDIR)/bin
DEP      = $(BIN)/deps

GO      = go
TIMEOUT = 15
V = 0
Q = $(if $(filter 1,$V),,@)
M = $(shell printf "\033[34;1m▶\033[0m")

export GO111MODULE=on

.PHONY: all
all: fmt vet lint test build ; $(info $(M) building executables…) @  ## Run all

$(DEP):
	@mkdir -p $@

$(DEP)/%: | $(DEP) ; $(info $(M) building $(PACKAGE)…)
	$Q tmp=$$(mktemp -d); \
	   env GO111MODULE=off GOPATH=$$tmp GOBIN=$(DEP) $(GO) get $(PACKAGE) \
		|| ret=$$?; \
	   rm -rf $$tmp ; exit $$ret

GOLINT = $(DEP)/golint
$(DEP)/golint: PACKAGE=golang.org/x/lint/golint

GOX = $(DEP)/gox
$(DEP)/gox: PACKAGE=github.com/mitchellh/gox

.PHONY: build
build: $(GOX) ; $(info $(M) building executables…) @  ## Build executables
	$Q $(GOX) -output "bin/{{.Dir}}_{{.OS}}_{{.Arch}}"

.PHONY: test
test: ; $(info $(M) running $(NAME:%=% )tests…) @ # Run tests
	$Q $(GO) test -v $(TESTPKGS)

.PHONY: lint
lint: | $(GOLINT) ; $(info $(M) running golint…) @ ## Run golint
	$Q $(GOLINT) -set_exit_status $(PKGS)

.PHONY: fmt
fmt: ; $(info $(M) running gofmt…) @ ## Run go fmt
	$Q $(GO) fmt $(PKGS)

.PHONY: vet
vet: ; $(info $(M) running go vet…) @ ## Run vet
	$Q $(GO) vet $(PKGS)

.PHONY: clean
clean: ; $(info $(M) cleaning…)	@ ## Clean up
	@rm -rf $(BIN)

.PHONY: help
help:
	@grep -hE '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-17s\033[0m %s\n", $$1, $$2}'