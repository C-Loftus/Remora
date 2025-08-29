build:
	wails build -tags webkit2_41 -v 2

dev:
	wails dev -tags webkit2_41 -v 2

tag:
	@read -p "Enter tag name: " tag; \
	read -p "Enter tag message: " msg; \
	git tag -a $$tag -m "$$msg"; \
	git push origin $$tag

.PHONY: build dev release