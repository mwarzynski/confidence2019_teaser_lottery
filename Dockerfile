FROM golang:1.12.0-alpine3.9

RUN mkdir -p /go/src/github.com/mwarzynski/confidence_web
ADD ./ /go/src/github.com/mwarzynski/confidence_web/
RUN cd /go/src/github.com/mwarzynski/confidence_web && go install

ENTRYPOINT confidence_web