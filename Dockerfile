FROM golang:alpine

RUN apk add --no-cache --update git

COPY . /go/src/github.com/nlnwa/sigridr-metrics
RUN cd /go/src/github.com/nlnwa/sigridr-metrics \
    && go get \
    && CGO_ENABLED=0 go build -a -tags netgo -v -ldflags "-w"

FROM scratch
LABEL maintainer="nettarkivet@nb.no"
COPY --from=0 /go/src/github.com/nlnwa/sigridr-metrics /
EXPOSE 8081
ENTRYPOINT ["/sigridr-metrics"]
