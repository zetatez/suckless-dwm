default: clean build

prepare:
	mkdir -p ./bins/
	go mod tidy

build-shell-cmds: prepare
	cp ./shell-cmds/* ./bins/

build: build-shell-cmds
	@cat ./go-build.sh
	./go-build.sh

clean:
	rm -rf bins/*

.PHONY: default
