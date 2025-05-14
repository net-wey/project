FROM golang:1.21 as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o task-calendar .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/task-calendar .
COPY --from=builder /app/static ./static
COPY --from=builder /app/migrations ./migrations

EXPOSE 80
CMD ["./task-calendar"]
