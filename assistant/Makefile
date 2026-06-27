.PHONY: build run dev install clean fmt swag

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o assistant ./cmd/assistant/

run: build
	@./assistant

dev:
	go run ./cmd/assistant --config config.yaml

install: build
	systemctl --user disable --now assistant 2>/dev/null || true
	sudo cp assistant /opt/assistant/assistant
	rm assistant
	rm -rf   ~/.config/assistant/scripts
	mkdir -p ~/.config/assistant/scripts
	cp -r scripts/* ~/.config/assistant/scripts/
	chmod +x ~/.config/assistant/scripts/*
	mkdir -p ~/.config/assistant/logs
	cp -f config.yaml ~/.config/assistant/config.yaml
	sed -i "s|filename:.*|filename: $$HOME/.config/assistant/logs/assistant.log|" ~/.config/assistant/config.yaml
	mkdir -p ~/.config/systemd/user
	sed "s|%h|$$HOME|g" assistant.service > ~/.config/systemd/user/assistant.service
	systemctl --user daemon-reload
	systemctl --user enable --now assistant 2>/dev/null || true
	@echo "installed"

clean:
	@rm -f assistant coverage.out

fmt:
	go fmt ./... && go vet ./...

swag:
	swag init -g ./cmd/assistant/main.go --output docs/
	@echo "Swagger Addr:"
	@echo "  http://127.0.0.1:4321/swagger/index.html"
