FROM golang:1.12.0-alpine3.9

RUN mkdir -p /go/src/github.com/mwarzynski/confidence2019_teaser_lottery
ADD ./ /go/src/github.com/mwarzynski/confidence2019_teaser_lottery/
RUN cd /go/src/github.com/mwarzynski/confidence2019_teaser_lottery && go install

ENTRYPOINT confidence2019_teaser_lottery