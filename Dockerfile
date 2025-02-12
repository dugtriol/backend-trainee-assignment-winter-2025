FROM golang:1.23.5

WORKDIR ${GOPATH}/avito-shop/
COPY . ${GOPATH}/avito-shop/

RUN go build -o /build ./cmd/app

EXPOSE 8080

CMD ["/build"]