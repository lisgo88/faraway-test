ENV=local
CACHE_TTL=10

CLIENT_HOST=localhost
CLIENT_PORT=8080
CLIENT_MAX_ATTEMPTS=10000

SERVER_HOST=localhost
SERVER_PORT=8080
SERVER_MAX_ATTEMPTS=10000
SERVER_HASH_BITS=3
SERVER_HASH_TTL=300
SERVER_TIMEOUT=30

test:
	@echo "Running tests..."
	go test ./... -cover -race

server-run:	server-build
	@echo "Running client..."
	./bin/server

server-build:
	@echo "Building server..."
	GOGC=off CGO_ENABLED=0 go build -v -o bin/server cmd/server/main.go

client-run:	client-build
	@echo "Running client..."
	./bin/client

client-build:
	@echo "Building client..."
	GOGC=off CGO_ENABLED=0 go build -v -o bin/client cmd/client/main.go