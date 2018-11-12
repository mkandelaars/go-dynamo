FROM golang:1.11-alpine

RUN apk add git
WORKDIR /go/src/app
COPY ./src/ .
COPY ./dist /go/src/

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["app"]
