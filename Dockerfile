FROM golang:1.18-alpine

WORKDIR /opt/catearsbot

RUN apk update \
 && apk add --no-cache \
            --virtual build \
            gcc \
            git \
            musl-dev

COPY . .

RUN go mod vendor \
 && go build main.go \
 && apk del build

CMD ["./main"]