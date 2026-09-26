FROM golang:1.27.1-alpine AS build

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 \
    go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /bin/betsystem \
    ./cmd/betsystem


FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=build /bin/betsystem /app/betsystem

EXPOSE 8080

USER 65532:65532

ENTRYPOINT ["/app/betsystem"]