FROM golang:1.11beta2-alpine
RUN apk add --update git gcc musl-dev
WORKDIR /usr/src

ADD go.mod go.sum ./
RUN go mod -sync
ADD . .
RUN CGO_ENABLED=0 go install ./...

FROM busybox
COPY --from=0 /go/bin/* /usr/local/bin/
ENTRYPOINT [ "renderd" ]
EXPOSE 8080
USER nobody
