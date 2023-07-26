APP_NAME=prompt-manager
VERSION?=$(shell git describe --tags)
GO_TARGET_OS?=darwin
GO_TARGET_ARCH?=arm64

custom_prompts.json:
	@echo "Creating custom prompts file..."
	@touch custom_prompts.json
	@echo "[]" > custom_prompts.json
	@echo "Custom prompts file created!"

.PHONY: env
env:
	@echo "Setting up environment variables..."
	@source env.sh
	@echo "Environment variables set up!"

.PHONY: build
build: custom_prompts.json
	@echo "Building..."
	GOOS=$(GO_TARGET_OS) GOARCH=$(GO_TARGET_ARCH) go build -o $(APP_NAME) src/main.go; 
	@echo "Build complete!"


.PHONY: archive 
archive: build
	@echo "Archiving..."
	@zip -r $(APP_NAME)_$(GO_TARGET_OS)_$(GO_TARGET_ARCH).alfredworkflow $(APP_NAME) custom_prompts.json icon.png info.plist
	@echo "Archive complete!"


.PHONY: version
version:
	@echo "{ \"version\": \"$(VERSION)\" }" > package.json
	@sed -i '' 's/<string>v[0-9]\.[0-9]\.[0-9]<\/string>/<string>$(VERSION)<\/string>/g' info.plist
