# syntax=docker/dockerfile:1.2

# Copyright 2020-2021 Authors of Cilium
# SPDX-License-Identifier: Apache-2.0

ARG BASE_IMAGE=scratch
ARG GOLANG_IMAGE=quay.io/cilium/cilium-builder:6615810f9d5595ea9639e91255bdbe83acb4b5fb@sha256:595c6a8c06db1dd2a71a4a2891b270af674705b7d1f7a2cfec5bfac64bd2af74
ARG ALPINE_IMAGE=docker.io/library/alpine:3.12.7@sha256:36553b10a4947067b9fbb7d532951066293a68eae893beba1d9235f7d11a20ad

# BUILDPLATFORM is an automatic platform ARG enabled by Docker BuildKit.
# Represents the plataform where the build is happening, do not mix with
# TARGETARCH
FROM --platform=${BUILDPLATFORM} ${GOLANG_IMAGE} as builder

# TARGETOS is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETOS
# TARGETARCH is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETARCH
ARG NOSTRIP

WORKDIR /go/src/github.com/cilium/tetragon
RUN --mount=type=bind,readwrite,target=/go/src/github.com/cilium/tetragon --mount=target=/root/.cache,type=cache --mount=target=/go/pkg/mod,type=cache \
    make GOARCH=${TARGETARCH} tetragon-operator-image \
    && mkdir -p /out/${TARGETOS}/${TARGETARCH}/usr/bin && mv tetragon-operator /out/${TARGETOS}/${TARGETARCH}/usr/bin

# BUILDPLATFORM is an automatic platform ARG enabled by Docker BuildKit.
# Represents the plataform where the build is happening, do not mix with
# TARGETARCH
FROM --platform=${BUILDPLATFORM} ${ALPINE_IMAGE} as certs
RUN apk --update add ca-certificates

FROM ${BASE_IMAGE}
# TARGETOS is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETOS
# TARGETARCH is an automatic platform ARG enabled by Docker BuildKit.
ARG TARGETARCH
LABEL maintainer="maintainer@cilium.io"
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /out/${TARGETOS}/${TARGETARCH}/usr/bin/tetragon-operator /usr/bin/tetragon-operator
WORKDIR /
ENV GOPS_CONFIG_DIR=/
CMD ["/usr/bin/tetragon-operator"]
