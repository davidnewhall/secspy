#
# building static go binary with Debian golang container
#

FROM golang:stretch as builder
ARG ARCH=amd64
ARG OS=linux
ARG BUILD_DATE=0
ARG COMMIT=0
ARG VERSION=development

# Build-time metadata as defined at http://label-schema.org
LABEL org.label-schema.build-date="${BUILD_DATE}" \
      org.label-schema.name="UniFi Poller" \
      org.label-schema.description="Polls a UniFi controller and stores metrics in InfluxDB" \
      org.label-schema.url="https://github.com/davidnewhall/unifi-poller" \
      org.label-schema.vcs-ref="${COMMIT}" \
      org.label-schema.vcs-url="https://github.com/davidnewhall/unifi-poller" \
      org.label-schema.vendor="Go Lift" \
      org.label-schema.architecture="${OS} ${ARCH}" \
      org.label-schema.license="MIT" \
      org.label-schema.version="${VERSION}" \
      org.label-schema.schema-version="1.0"

RUN mkdir -p $GOPATH/pkg/mod $GOPATH/bin $GOPATH/src/github.com/davidnewhall/secspy
RUN apt-get update \
  && apt-get install -y curl  \
  && curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

COPY . $GOPATH/src/github.com/davidnewhall/secspy
WORKDIR $GOPATH/src/github.com/davidnewhall/secspy

RUN dep ensure --vendor-only \
  && CGO_ENABLED=0 make secspy.${ARCH}.${OS} \
  && mv secspy.${ARCH}.${OS} secspy

FROM scratch

COPY --from=builder /go/src/github.com/davidnewhall/secspy/secspy /secspy
COPY --from=builder /go/src/github.com/davidnewhall/secspy/examples/secspy.conf.example /etc/secspy/up.conf

VOLUME [ "/etc/secspy" ]
ENTRYPOINT [ "/secspy" ]
