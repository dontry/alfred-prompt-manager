APP_NAME=prompt-manager
VERSION?=$(shell git describe --tags)
GO_TARGET_OS?=darwin
GO_TARGET_ARCH?=arm64
OUTPUT_FILE := $(APP_NAME)_$(GO_TARGET_OS)_$(GO_TARGET_ARCH)



custom_prompts.json:
	@echo "Creating custom prompts file..."
	@touch custom_prompts.json
	@echo "[]" > custom_prompts.json
	@echo "Custom prompts file created!"


.PHONY: build
build:  custom_prompts.json
	@echo "Building..."
	@GOOS=$(GO_TARGET_OS) GOARCH=$(GO_TARGET_ARCH) go build -o "$(OUTPUT_FILE)" src/main.go
	@echo "Build complete!"
	@go env > xx.txt


.PHONY: archive 
archive: build
	@echo "Archiving..."
	@zip -r $(OUTPUT_FILE).alfredworkflow $(OUTPUT_FILE) custom_prompts.json icon.png info.plist
	@echo "Archive complete!"


.PHONY: version
version:
	@echo "{ \"version\": \"$(VERSION)\" }" > package.json
	@sed -i '' 's/<string>v[0-9]\.[0-9]\.[0-9]<\/string>/<string>$(VERSION)<\/string>/g' info.plist
