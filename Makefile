SHELL := bash
.ONESHELL:
.SHELLFLAGS := -eu -o pipefail -c
.DELETE_ON_ERROR:
MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

ifeq ($(origin .RECIPEPREFIX), undefined)
  $(error This Make does not support .RECIPEPREFIX. Please use GNU Make 4.0 or later)
endif
.RECIPEPREFIX = >

build: out/aoc
.PHONY: build

install: tmp/cli.installed
.PHONY: install


tmp/cli.installed: out/aoc
> mkdir -p $(@D)
> mv out/aoc /usr/local/bin/
> touch $@

out/aoc: $(shell find ./ -type f  -regex ".*\.go")
> mkdir -p $(@D)
> go build -o out/aoc
> chmod +x out/aoc