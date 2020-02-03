# STEP 1 build executable binary
FROM golang:alpine as builder
MAINTAINER Statful Developer <developer@statful.com>

# Install SSL ca certificates
RUN apk update && \
    apk add git ca-certificates

# Create application user
RUN adduser -D -g '' app

WORKDIR /build

# cache go modules
COPY go.mod .
RUN go mod download

#build the binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s"

# STEP 2 build a small image
# start from scratch
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd

# Copy the static executable
COPY --from=builder /build/prometheus-exporter /prometheus-exporter
USER app

ENTRYPOINT ["/prometheus-exporter"]
