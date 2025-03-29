.PHONY: swag
swag:
	@swag init --dir ./src/main --output ./docs --generatedTime --parseInternal --parseDependency --parseDepth 3
