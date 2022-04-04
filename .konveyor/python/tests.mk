### Important, the LANG_DIR must match one from Makefile
LANG_DIR=python

$(LANG_DIR)-tests:
	$(MAKE) pre-commit

pre-commit: setup-pre-commit
	. ${PYTHON_VENV}/bin/activate && pre-commit

isort: setup-pre-commit
	. ${PYTHON_VENV}/bin/activate && pre-commit run isort

flake8: setup-pre-commit
	. ${PYTHON_VENV}/bin/activate && pre-commit run flake8
