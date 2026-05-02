FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY backend ./backend
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/imacan-api ./backend/cmd/api

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=build /out/imacan-api /app/imacan-api

EXPOSE 8080

CMD ["/app/imacan-api"]
