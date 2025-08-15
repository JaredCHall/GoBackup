# ---- build stage ----
FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# static-ish build is fine here
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gobackup .

# ---- runtime stage ----
FROM alpine:3.20
RUN apk add --no-cache bash rsync ca-certificates tzdata tree tini
WORKDIR /app
COPY --from=builder /gobackup /usr/local/bin/gobackup
# Keep sandbox running and respect SIGTERM/SIGINT
ENTRYPOINT ["/sbin/tini","--"]
CMD ["sleep", "infinity"]