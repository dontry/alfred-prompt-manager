
APP_NAME=alfred-prompt-manager

.PHONY: env
env:
	@echo "Setting up environment variables..."
	@sh env.sh
	@echo "Environment variables set up!"

.PHONY: build
build: env
	@echo "Building..."
	@go build -o bin/$(APP_NAME) src/main.go
	@echo "Build complete!"