.PHONY: build clean dev-api dev-web lint test cover

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

lint:
	golangci-lint run ./...
	cd web && npx eslint .
	cd web && npx prettier --check .

test:
	go test ./...
	cd web && npx vitest run

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -func=coverage.out
	cd web && npx vitest run --coverage
