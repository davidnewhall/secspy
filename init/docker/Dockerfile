#
# building static go binary with Debian golang container
#

FROM golang:stretch as builder
ARG ARCH=amd64
ARG OS=linux

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
