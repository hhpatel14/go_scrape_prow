### Important, the LANG_DIR must match one from Makefile
LANG_DIR=python
 
PYTHON_BINARY=python3

ifeq (, $(shell which $(PYTHON_BINARY) ))
  $(error "PYTHON=$(PYTHON_BINARY) binary not found in $(PATH)")
endif

ifndef PYTHON_VENV
  PYTHON_VENV=.venv
endif

$(PYTHON_VENV):
	test -d ${PYTHON_VENV} || ${PYTHON_BINARY} -m venv ${PYTHON_VENV}
	. ${PYTHON_VENV}/bin/activate && pip install -U pip
	touch ${PYTHON_VENV}

$(LANG_DIR)-dev-env: $(PYTHON_VENV)
	$(info **** To run VENV: $$source ${PYTHON_VENV}/bin/activate)
	$(info **** To later deactivate VENV: $$deactivate)
	$(info **** To destroy: make clean-dev-env)

$(LANG_DIR)-clean-dev-env:
	rm -rf ${PYTHON_VENV}
	find . -iname "*.pyc" -delete
