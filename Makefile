TARGETS := $(shell ls scripts | grep -vE 'clean|dev|dockerlint|help|release|shellcheck|tag')

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

.PHONY: tag
tag:
	./scripts/tag

.PHONY: shellcheck
shellcheck:
	@for file in $(shell find . -type f -executable -not -path "./.git/*" -not -path "./vendor/*"); do \
		echo "Validating : $$file"; \
		docker container run --rm --mount type=bind,src=$$PWD,dst=/mnt,ro koalaman/shellcheck -e SC2086 -e SC2046 -e SC1090 "$$file"; \
		if [ $$? -gt 0 ]; then \
			continue; \
		fi; \
	done;

.PHONY: dockerlint
dockerlint:
	@for file in $(shell find . -name 'Dockerfile*'); do \
		echo "Validating : $$file"; \
		docker container run -i --rm hadolint/hadolint hadolint --ignore DL3018 --ignore DL3013 - < $$file; \
	done;

.DEFAULT_GOAL := ci
