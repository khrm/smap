FROM registry.bookmyshow.org/golang:1.9.2 as build-img
LABEL maintainer "khrm.baig@gmail.com"

RUN mkdir -p /go/src/github.com/khrm/smap
WORKDIR /go/src/github.com/khrm/smap
ADD ./ /go/src/github.com/khrm/smap/

RUN CGO_ENABLED=0 go build -race -a -ldflags "-s -w" -o /smap cmd/smap/*.go


FROM alpine:latest


RUN apk --no-cache add ca-certificates
WORKDIR /usr/local/bin
COPY --from=build-img /smap .

CMD smap
