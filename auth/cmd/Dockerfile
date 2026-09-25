FROM golang:1.25.5-bookworm AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/exe ./cmd 

FROM alpine:3.23 AS runner
WORKDIR /app
COPY --from=builder /app/exe /app/
EXPOSE 8081
CMD ["/app/exe"]