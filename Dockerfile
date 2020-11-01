FROM golang:1.15 as builder

COPY / /polldance
WORKDIR /polldance

RUN make -e CGO_ENABLED=0 polldance

FROM alpine:3
RUN apk add --update-cache \
    jq \
  && rm -rf /var/cache/apk/*
COPY --from=builder /polldance/polldance /usr/bin
