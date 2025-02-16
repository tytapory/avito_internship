FROM golang:1.23

WORKDIR /app
COPY . /app/

RUN go build -o /build ./cmd/app \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]