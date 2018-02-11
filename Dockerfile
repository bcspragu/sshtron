FROM golang:1.9.4-alpine3.7 as builder

WORKDIR $GOPATH/src/github.com/bcspragu/sshtron
ADD . .
RUN apk update && apk add git && go get && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /usr/bin/sshtron .

FROM alpine:latest

COPY --from=builder /usr/bin/sshtron /usr/bin/
RUN apk add --update openssh-client && ssh-keygen -t rsa -N "" -f id_rsa
ENTRYPOINT sshtron
