SHELL = bash

.PHONY: default goimports lint vet help clean build

default:  goimports lint vet build ## Run default target : fmt + lints + build

goimports:  ## Run goimports to format code
	goimports -w .

lint:  ## Lint all go code in project
	golint ./...

vet:  ## Go vet all project code
	go vet ./...

build: clean ## Go build
	go build -o bin/fetchdou cmd/fetchdou/*go

help:  ## Show This Help
	@for line in $$(cat Makefile | grep "##" | grep -v "grep" | sed  "s/:.*##/:/g" | sed "s/\ /!/g"); do verb=$$(echo $$line | cut -d ":" -f 1); desc=$$(echo $$line | cut -d ":" -f 2 | sed "s/!/\ /g"); printf "%-30s--%s\n" "$$verb" "$$desc"; done

clean:  ## Clean up transient (generated) files
	go clean
	rm -f bin/*
