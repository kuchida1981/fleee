.PHONY: build clean dev-api dev-web

build:
	cd web && npm run build
	go build -o fleee ./cmd/fleee

clean:
	rm -f fleee
	rm -rf web/dist

dev-api:
	go run ./cmd/fleee serve

dev-web:
	cd web && npm run dev
