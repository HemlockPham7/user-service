### Base ###

FROM golang:1.26-alpine AS base

RUN mkdir -p /opt/app

WORKDIR /opt/app

COPY . .

RUN go mod download

RUN apk add build-base

### Build ###

FROM base AS build

RUN GOOS=linux go build -tags musl -ldflags "-w -s" \
    -o bookmark-service cmd/api/main.go

### TEST-EXEC ###

FROM base AS test-exec

ARG _outputdir="/tmp/coverage"
ARG COVERAGE_EXCLUDE

RUN mkdir -p ${_outputdir} && \
    CGO_ENABLED=1 go test ./... -coverprofile=coverage.tmp -coverpkg=./... -covermode=atomic -p 1 && \
    grep -vE "${COVERAGE_EXCLUDE}" coverage.tmp > ${_outputdir}/coverage.out && \
    go tool cover -html=${_outputdir}/coverage.out -o ${_outputdir}/coverage.html

### Test ###

FROM scratch AS test

ARG _outputdir="/tmp/coverage"
COPY --from=test-exec ${_outputdir}/coverage.out /
COPY --from=test-exec ${_outputdir}/coverage.html /

### Final ###

FROM alpine:3.24.1 AS final

WORKDIR /app

COPY --from=build /opt/app/bookmark-service /app/bookmark-service
COPY --from=build /opt/app/docs /app/docs

CMD ["/app/bookmark-service"]
