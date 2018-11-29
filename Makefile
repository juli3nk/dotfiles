TARGETS := $(shell ls scripts | grep -vE 'clean|dev|help|release')

.PHONY: .dapper
.dapper:
	@echo Downloading dapper
	@curl -sL https://releases.rancher.com/dapper/latest/dapper-`uname -s`-`uname -m|sed 's/v7l//'` > .dapper.tmp
	@@chmod +x .dapper.tmp
	@./.dapper.tmp -v
	@mv .dapper.tmp .dapper

.PHONY: .github-release
.github-release:
	@echo Downloading github-release
	@curl -sL https://github.com/aktau/github-release/releases/download/v0.7.2/linux-amd64-github-release.tar.bz2 | tar xjO > .github-release.tmp
	@@chmod +x .github-release.tmp
	@./.github-release.tmp -v
	@mv .github-release.tmp .github-release

.PHONY: $(TARGETS)
$(TARGETS): .dapper
	./.dapper $@

.PHONY: clean
clean:
	@./scripts/clean

.PHONY: dev
dev: .dapper
	./.dapper -f Dockerfile.dev -m bind -s

.PHONY: help
help:
	@./scripts/help

.PHONY: release
release: .github-release
	./scripts/release

.DEFAULT_GOAL := ci
