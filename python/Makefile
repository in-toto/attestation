SHELL := /bin/bash

PY_MODULE := in_toto_attestation

ALL_PY_SRCS := $(shell find $(PY_MODULE) tests -name '*.py')

# Optionally overriden by the user, if they're using a virtual environment manager.
VENV ?= env

VENV_STAMP = $(VENV)/pyvenv.cfg

# On Windows, venv scripts/shims are under `Scripts` instead of `bin`.
VENV_BIN := $(VENV)/bin
ifeq ($(OS),Windows_NT)
	VENV_BIN := $(VENV)/Scripts
endif

# Optionally overridden by the user in the `release` target.
BUMP_ARGS :=

# Optionally overridden by the user/CI, to limit the installation to a specific
# subset of development dependencies.
INSTALL_EXTRA := dev

.PHONY: all
all:
	@echo "Run my targets individually!"

.PHONY: dev
dev: $(VENV_STAMP)

$(VENV_STAMP): pyproject.toml
	python -m venv $(VENV)
	$(VENV_BIN)/python -m pip install --upgrade pip
	$(VENV_BIN)/python -m pip install -e .[$(INSTALL_EXTRA)]

.PHONY: lint
lint: $(VENV_STAMP)
	. $(VENV_BIN)/activate && \
		ruff format --check $(ALL_PY_SRCS) && \
		ruff $(ALL_PY_SRCS) && \
		mypy $(PY_MODULE)

.PHONY: reformat
reformat: $(VENV_STAMP)
	. $(VENV_BIN)/activate && \
		ruff format $(ALL_PY_SRCS) && \
		ruff --fix $(ALL_PY_SRCS)

.PHONY: test tests
test tests: $(VENV_STAMP)
	. $(VENV_BIN)/activate && \
		pytest --cov=$(PY_MODULE) $(T) $(TEST_ARGS) && \
		python -m coverage report -m $(COV_ARGS)

.PHONY: package
package: $(VENV_STAMP)
	. $(VENV_BIN)/activate && \
		python -m build

.PHONY: edit
edit:
	$(EDITOR) $(ALL_PY_SRCS)

.PHONY: clean
clean:
	rm -r $(VENV)
