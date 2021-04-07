FROM golang:alpine AS builder

ADD ./ /go/src/github.com/ravens/modem_exporter/

RUN set -ex && \
  cd /go/src/github.com/ravens/modem_exporter && \       
  CGO_ENABLED=0 go build \
        -v -a \
        -ldflags '-extldflags "-static"' && \
  mv ./modem_exporter /usr/bin/modem_exporter

FROM busybox
COPY --from=builder /usr/bin/modem_exporter /usr/local/bin/modem_exporter
ENTRYPOINT [ "modem_exporter" ]