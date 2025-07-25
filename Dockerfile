FROM golang:alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o healthcheck .

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/healthcheck .
EXPOSE 8081
CMD ["./healthcheck"]