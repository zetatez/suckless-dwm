# Assistant

Local service that exposes system commands and utilities via HTTP API. Triggered from dwm keybindings or any HTTP client.

## Quick Start

```bash
# One-line install
curl -sL https://github.com/zetatez/suckless-dwm/raw/master/assistant/install.sh | sh

# Or manually:
# 1. make install
# 2. systemctl --user enable --now assistant
```

## API

All endpoints under `http://<host>:4321/api/`.
Auth: HTTP Basic (`auth.username` / `auth.password` from config).

See `scripts/` for all available endpoints. Each file is a single curl command.

## Scripts

Each script in `scripts/` is a single curl command for one API endpoint:

```bash
./scripts/sys-shortcut               # rofi 电源菜单
./scripts/sys-display                # rofi 显示器布局
./scripts/toggle.sh flameshot
./scripts/launch.sh inkscape
./scripts/open-url-chrome.sh "https://github.com"
```

Install with `make install` or `systemctl --user enable --now assistant`.

## Project Structure

```
cmd/assistant/        # entrypoint
internal/
├── app/
│   ├── server.go     # Gin server setup
│   └── modules/
│       ├── health/       # health check
│       ├── svc/          # all API endpoints (handler + service)
│       └── filebrowser/  # /api/files/* (list/raw/download/upload + embed UI)
├── bootstrap/        # init config, log, shutdown
pkg/
├── response/         # API response helpers
├── xlog/             # logging
└── llm/              # LLM abstraction (unused, kept for future)
```
