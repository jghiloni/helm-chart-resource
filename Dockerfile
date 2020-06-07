FROM golang:1 AS builder

COPY . /src
WORKDIR /src

ENV CGO_ENABLED 0

RUN go get -d ./...
RUN go build -o /assets/check ./cmd/check
RUN go build -o /assets/in ./cmd/in
RUN set -e; for pkg in $(go list ./...); do \
        go test -o "/tests/$(basename $pkg).test" -c $pkg; \
    done

FROM alpine:edge AS resource
RUN apk add --no-cache bash tzdata ca-certificates
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
COPY --from=builder /tests /tests
ADD . /docker-image-resource
RUN set -e; for test in /tests/*.test; do \
        $test; \
    done

FROM resource
