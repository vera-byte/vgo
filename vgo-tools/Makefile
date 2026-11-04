.PHONY: help
help: ## 查看帮助
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'



.PHONY: pack.template-simple
pack.template-simple:
	@rm -fr temp
	@mkdir temp || exit 0
	@cd temp && git clone --depth=1 -b simple  https://github.com/vera-byte/vgo.git vgo-simple
	@rm -fr temp/vgo-simple/.git
	@cd temp && gf pack vgo-simple ../internal/packed/vgo-simple.go -n=packed -y
	@rm -fr temp/vgo-simple

.PHONY: pack.template-simple.ssh
pack.template-simple.ssh:
	@rm -fr temp
	@mkdir temp || exit 0
	@cd temp && git clone --depth=1 -b simple  git@github.com:vera-byte/vgo.git vgo-simple
	@rm -fr temp/vgo-simple/.git
	@cd temp && gf pack vgo-simple ../internal/packed/vgo-simple.go -n=packed -y
	@rm -fr temp

.PHONY: pack.docs
pack.docs: ## 导出 Mintlify 静态站并打入 vgo-tools（发布期执行；需 Mintlify CLI：npm i -g mint）
	@set -e; rm -fr temp; mkdir temp
	@cd ../docs && mint export --output ../vgo-tools/temp/export.zip
	@cd temp && mkdir docs && cd docs && unzip -q ../export.zip && cd .. && rm -f export.zip
	@cd temp && gf pack docs ../internal/packed/docs.go -n=packed -y
	@rm -fr temp

# Install/Update to the latest CLI tool.
.PHONY: cli
cli:
	@set -e; \
	wget -O gf https://github.com/gogf/gf/releases/latest/download/gf_$(shell go env GOOS)_$(shell go env GOARCH) && \
	chmod +x gf && \
	./gf install  && \
	rm ./gf



# Install/Update to the latest vgo-tools.
.PHONY: tools
tools:
	@set -e; \
	curl -L https://download.fastgit.org/vera-byte/vgo/releases/latest/download/vgo-tools_$(shell go env GOOS)_$(shell go env GOARCH) -o ./vgo-tools && \
	chmod +x vgo-tools && \
	./vgo-tools install && \
	rm ./vgo-tools


# Check and install vgo-tools.
.PHONY: tools.install
tools.install:
	@set -e; \
	echo "Checking vgo-tools..."; \
	vgo-tools -v > /dev/null 2>&1 || if [[ "$?" -ne "0" ]]; then \
  		echo "vgo-tools is not installed, start proceeding auto installation..."; \
		make tools; \
	fi;\
	echo "vgo-tools is installed.";