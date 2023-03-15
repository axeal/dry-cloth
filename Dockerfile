# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS builder

WORKDIR $GOPATH/src/axeal/dry-cloth/
COPY . .

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /go/bin/dry-cloth cmd/dry-cloth/main.go

FROM gcr.io/distroless/static AS final

USER nonroot:nonroot
COPY --from=builder --chown=nonroot:nonroot /go/bin/dry-cloth /go/bin/dry-cloth
ENTRYPOINT ["/go/bin/dry-cloth"]
