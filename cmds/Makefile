default: clean build

prepare:
	mkdir -p ./bins/

build-shell-cmds: prepare
	cp -f ./shell-cmds/* ./bins/

build: build-shell-cmds
	@cat ./go-build.sh
	go mod tidy
	./go-build.sh

clean:
	rm -rf bins/*

.PHONY: default
