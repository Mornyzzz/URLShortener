FROM golang:1.23-alpine AS builder

WORKDIR /usr/local/src
RUN apk --no-cache add bash git make gcc gettext musl-dev

COPY ["app/go.mod", "app/go.sum", "./"]
RUN go mod download


COPY app ./
RUN go build -o ./bin/app cmd/url_shortener/main.go


FROM alpine:latest

COPY --from=builder /usr/local/src/bin/app /
COPY config/config.yaml /config/config.yaml

EXPOSE 8080

CMD ["./app", "${STORAGE_TYPE}"]
