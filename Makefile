GIT_TAG := $(shell git describe --tags --exact-match --abbrev=0 2>/dev/null)
BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
IMG_NAME := hemlockpham/user-service
IMG_TAG := latest

ifneq ($(GIT_TAG),)
	IMG_TAG := $(GIT_TAG)
endif

export IMG_TAG

COVERAGE_EXCLUDE=mocks|main.go|test|config.go|client.go|level.go|mock.go|request.go|parser.go|postgres_data|infrastructure|migration.go
COVERAGE_THRESHOLD=70
COVERAGE_FOLDER=./test-output

docker-test:
	mkdir -p $(COVERAGE_FOLDER)
	docker buildx build --build-arg COVERAGE_EXCLUDE="${COVERAGE_EXCLUDE}" --progress=plain --target test -t test:test --output $(COVERAGE_FOLDER) .
	@total=$$(go tool cover -func=$(COVERAGE_FOLDER)/coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
	   echo "❌ Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
	   exit 1; \
	else \
	   echo "✅ Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

docker-build:
	docker build -t $(IMG_NAME):$(IMG_TAG) .

DOCKER_USERNAME ?=
DOCKER_PASSWORD ?=

docker-login:
	echo "$(DOCKER_PASSWORD)" | docker login -u "$(DOCKER_USERNAME)" --password-stdin

docker-release:
	docker push $(IMG_NAME):$(IMG_TAG)

run:
	go run cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go --output docs

dev-run: swagger run