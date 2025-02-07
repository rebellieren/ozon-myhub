
FROM golang:1.23.6-alpine3.20 AS builder


RUN apk add --no-cache tzdata


ENV CGO_ENABLED=0


WORKDIR /app


COPY go.mod go.sum ./


RUN go mod download


COPY . .


RUN go build -ldflags="-s -w" -o /app/bin/myhub-app ./cmd/app


FROM alpine:3.20


WORKDIR /app


COPY --from=builder /app/bin/myhub-app .


ENTRYPOINT [ "./myhub-app" ]
