go := $(shell which go)
watch_files := (.go$$)|(.html$$)|(.js$$)

setup: install
install:
	$(go) get -u golang.org/x/tools/cmd/goimports
	$(go) get -u github.com/cespare/reflex
	$(go) mod tidy

run:
	$(go) run cmd/example/main.go

watch:
	reflex -r "$(watch_files)" -s -- make run
