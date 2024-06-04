FROM golang:1.21.6-alpine

WORKDIR /app

RUN apk add build-base

COPY . /app

RUN go mod download

RUN go build -o app .

EXPOSE ${PORT}

CMD ./app