FROM golang:1.23 AS yq

RUN set -x \
  && apt-get update \
  && apt-get install -y --no-install-recommends apt-transport-https ca-certificates

## build yq binary from ac-3.3.4 branch which has updated go dependencies
RUN set -x \
  && git clone https://github.com/gomodules/yq.git \
  && cd yq \
  && git checkout ac-3.3.4 \
  && CGO_ENABLED=0 go build -v -o /yq yq.go

FROM alpine:3.20.2

RUN apk add --no-cache bash

COPY scripts/config-merger.sh /usr/local/bin/config-merger.sh
COPY --from=yq /yq /usr/bin/yq
RUN ["mkdir","-p","/data"]

RUN chmod -c 755 /usr/local/bin/config-merger.sh
#RUN chown -R 65534:65534 /data

#USER 65534

ENTRYPOINT ["/bin/bash","/usr/local/bin/config-merger.sh"]
