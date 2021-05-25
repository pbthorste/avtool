# Build Variables
VERSION ?= $(shell git describe --tags --always)

# Go variables
GO      ?= go
GOOS    ?= $(shell $(GO) env GOOS)
GOARCH  ?= $(shell $(GO) env GOARCH)
GOHOST  ?= GOOS=$(GOOS) GOARCH=$(GOARCH) $(GO)

LDFLAGS ?= "-X main.version=$(VERSION)"

.PHONY: all
all: help

###############
##@ Development

.PHONY: clean
clean: ## Clean workspace
	@ $(MAKE) --no-print-directory log-$@
	rm -rf coverage.out
	go mod tidy

.PHONY: test
test: ## Run tests
	@ $(MAKE) --no-print-directory log-$@
	$(GOHOST) test -covermode count -coverprofile coverage.out -v  -run ^Test ./...

.PHONY: lint
lint: ## Run linters
	@ $(MAKE) --no-print-directory log-$@
	golangci-lint run

###########
##@ Release

.PHONY: check
check:  ## Check if version exists
	@ $(MAKE) --no-print-directory log-$@
 	ifeq ($(shell git tag -l | grep -c $(VERSION) 2>/dev/null), 0)
 	else
		@git tag -l
		@echo "VERSION: $(VERSION) has already been used. Please pass in a different version value."
		@exit 2
 	endif

.PHONY: changelog
changelog: check  ## Generate changelog
	@ $(MAKE) --no-print-directory log-$@
	git-chglog --next-tag $(VERSION) -o CHANGELOG.md

.PHONY: release
release: changelog   ## Release a new tag
	@ $(MAKE) --no-print-directory log-$@
	git add CHANGELOG.md
	git commit -m "chore: update changelog for $(VERSION)"
	git tag $(VERSION)
	git push origin main $(VERSION)

########
##@ Help

.PHONY: help
help: ## Display this help
	@awk \
		-v "col=\033[36m" -v "nocol=\033[0m" \
		' \
			BEGIN { \
				FS = ":.*##" ; \
				printf "Usage:\n  make %s<target>%s\n", col, nocol \
			} \
			/^[a-zA-Z_-]+:.*?##/ { \
				printf "  %s%-12s%s %s\n", col, $$1, nocol, $$2 \
			} \
			/^##@/ { \
				printf "\n%s%s%s\n", nocol, substr($$0, 5), nocol \
			} \
		' $(MAKEFILE_LIST)

log-%:
	@grep -h -E '^$*:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk \
			'BEGIN { \
				FS = ":.*?## " \
			}; \
			{ \
				printf "\033[36m==> %s\033[0m\n", $$2 \
			}'
