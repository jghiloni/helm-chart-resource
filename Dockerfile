FROM golang:1 AS builder

COPY . /src
WORKDIR /src

ENV CGO_ENABLED 0

RUN go get -d ./...
RUN go build -o /assets/check ./cmd/check
RUN go build -o /assets/in ./cmd/in
RUN cp ./cmd/out/out.sh /assets/out
RUN set -e; for pkg in $(go list ./...); do \
        go test -o "/tests/$(basename $pkg).test" -c $pkg; \
    done

FROM alpine:edge AS resource
ARG VERSION=3.2.2

RUN apk add --update --no-cache bash jq git tzdata curl ca-certificates && \
    curl -L https://get.helm.sh/helm-v${VERSION}-linux-amd64.tar.gz |tar xvz && \
    mv linux-amd64/helm /usr/bin/helm && \
    chmod +x /usr/bin/helm && \
    rm -rf linux-amd64 && \
    apk del curl && \
    rm -f /var/cache/apk/*
RUN helm plugin install https://github.com/chartmuseum/helm-push.git
COPY --from=builder assets/ /opt/resource/
RUN chmod +x /opt/resource/*

FROM resource AS tests
COPY --from=builder /tests /tests
ADD . /docker-image-resource
RUN set -e; for test in /tests/*.test; do \
        $test; \
    done

FROM resource
