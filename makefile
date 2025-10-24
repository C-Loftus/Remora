build:
	wails build -tags webkit2_41 -v 2

dev:
	wails dev -tags webkit2_41 -v 2

release:
	git add . 
	git commit -m "release" || true
	git push origin main
	@read -p "Enter tag name: " tag; \
	read -p "Enter tag message: " msg; \
	git tag -a $$tag -m "$$msg"; \
	git push origin $$tag

# install the binary and its associated .desktop file
# so it shows up in the applications menu
install:
	cp ./remora.desktop ~/.local/share/applications/remora.desktop
	cp ./frontend/src/assets/images/remora.png ~/.local/share/icons/remora.png
	wails build -tags webkit2_41 -v 2
	cp build/bin/remora  ~/.local/bin/remora
	chmod +x ~/.local/bin/remora
	update-desktop-database ~/.local/share/applications/

.PHONY: build dev release install