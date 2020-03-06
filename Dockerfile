FROM golang:alpine AS builder

RUN mkdir -p /go/src/binlog-db-sync
WORKDIR /go/src/binlog-db-sync
COPY . /go/src/binlog-db-sync

RUN apk --no-cache add git

RUN go get -u github.com/aws/aws-sdk-go/aws
RUN go get -u github.com/aws/aws-sdk-go/aws/session
RUN go get -u github.com/aws/aws-sdk-go/service/lambda
RUN go get -u github.com/json-iterator/go
RUN go get -u github.com/siddontang/go-mysql/canal
RUN go get -u github.com/siddontang/go-mysql/schema
RUN go get -u gopkg.in/yaml.v2

RUN go build -o binlog-db-sync-server

FROM alpine

RUN mkdir -p /usr/src/app
WORKDIR /usr/src/app

COPY --from=builder /go/src/binlog-db-sync/binlog-db-sync-server .
COPY ./handlers-settings.yaml .

CMD ["./binlog-db-sync-server"]