
APP_NAME=alfred-prompt-manager

custom_prompts.json:
	@echo "Creating custom prompts file..."
	@touch custom_prompts.json
	@echo "[]" > custom_prompts.json
	@echo "Custom prompts file created!"

.PHONY: env
env:
	@echo "Setting up environment variables..."
	@sh env.sh
	@echo "Environment variables set up!"

.PHONY: build
build: env custom_prompts.json
	@echo "Building..."
	@go build src/main.go
	@echo "Build complete!"