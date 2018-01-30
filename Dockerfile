FROM golang:1.9.2-alpine

RUN apk add --no-cache --update alpine-sdk protobuf protobuf-dev

COPY . /go/src/git.nb.no/nettarkiv/sigridr-metrics
RUN cd /go/src/git.nb.no/nettarkiv/sigridr-metrics && make release-binary

FROM scratch
LABEL maintainer="nettarkivet@nb.no"
COPY --from=0 /go/src/git.nb.no/nettarkiv/sigridr-metrics /
EXPOSE 8081
ENTRYPOINT ["/sigridr-metrics"]
