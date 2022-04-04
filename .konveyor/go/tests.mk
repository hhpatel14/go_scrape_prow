### Important, the LANG_DIR must match one from Makefile
LANG_DIR=go

$(LANG_DIR)-tests:
	$(warning I AM HERE)
	go build . && \
	go test -v 

go-lint: setup-pre-commit
	. ${PYTHON_VENV}/bin/activate && pre-commit run golangci-lint
