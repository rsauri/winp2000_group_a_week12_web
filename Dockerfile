FROM golang:1.23 AS builder
WORKDIR /web
COPY . .
#RUN CSG_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-app .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go-app .

#FROM alpine:3.20
FROM ubuntu:22.04

# Install dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    ca-certificates \
    curl \
    && rm -rf /var/lib/apt/lists/*
    
WORKDIR /srv
COPY --from=builder  /web/go-app .
COPY --from=builder  /web/site ./site

EXPOSE 80
CMD [ "./go-app" ]