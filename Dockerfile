FROM golang:1.11beta2-alpine
RUN apk add --update git gcc musl-dev
WORKDIR /usr/src

ADD . .
RUN go install ./...

FROM busybox
COPY --from=0 /go/bin/* /usr/local/bin/
ENTRYPOINT [ "renderd" ]
