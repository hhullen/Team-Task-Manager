FROM golang:1.25-alpine3.23 AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build ./cmd/migrator

FROM alpine:3.23.0

WORKDIR /app
COPY --from=builder /build .

ENV RUNNING_IN_CONTAINER=true

RUN chmod +x ./migrator

ENTRYPOINT [ "./migrator", "up" ]