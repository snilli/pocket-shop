BUILD_TAGS = -tags=sonic

.PHONY: run-api build-api swagger mock test test-integration test-coverage test-watch clean

run-api:
	go run $(BUILD_TAGS) cmd/api/main.go

build-api:
	mkdir -p bin
	go build $(BUILD_TAGS) -o bin/api cmd/api/main.go

swagger:
	swag init -g internal/infrastructure/server/server.go -o docs

mock:
	mockery

test:
	ginkgo run --randomize-all --cover -r --label-filter="!integration"

test-integration:
	ginkgo run --randomize-all -r --label-filter="integration"

test-coverage:
	ginkgo run --randomize-all --cover -r --coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

test-watch:
	ginkgo watch -r

clean:
	rm -rf bin/ coverage.out coverage.html

