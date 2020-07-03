FROM golang:1 as builder

COPY . /src

WORKDIR /src

ARG VERSION=dev

# https://github.com/golang/go/wiki/Modules#how-do-i-use-vendoring-with-modules-is-vendoring-going-away
# go build -mod=vendor
# The -gcflags "all=-N -l" flag helps us get a better debug experience
RUN set -x \
    && export BUILD=$(date -u +%s%N) \
    && export REPO=$(go list -m) \
    && export LDFLAGS="-X ${REPO}/misc.Version=${VERSION} -X ${REPO}/misc.Build=${BUILD}" \
    && export GOGC=off \
    && export CGO_ENABLED=0 \
    && [ -d "./vendor" ] || go mod vendor \
    && go build -v -mod=vendor -trimpath -gcflags "all=-N -l" -ldflags "${LDFLAGS}" -o /go/bin/neofs-s3 ./main.go

# Executable image
FROM alpine:3.10

WORKDIR /

COPY --from=builder /go/bin/neofs-s3 /usr/bin/neofs-s3
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Run delve
CMD ["/usr/bin/neofs-s3"]
