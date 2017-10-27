ARG GOLANG_TAG=latest
ARG ALPINE_TAG=latest

# Build stage with all the development dependencies
FROM golang:${GOLANG_TAG} AS dev

# Setup gin for dev monitoring
RUN go-wrapper download github.com/codegangsta/gin
RUN go-wrapper install github.com/codegangsta/gin

# Copy the local package files into the image's workspace.
COPY . /go/src/github.com/datatogether/coverage
WORKDIR /go/src/github.com/datatogether/coverage

# Run tests
RUN go test

# Build the static api binary for production
RUN CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo

# Set binary as the default command
CMD ["coverage"]

# Start over from an Alpine Linux image as a base
# to create a minumal production image
FROM alpine:${ALPINE_TAG}
LABEL repo="https://github.com/datatogether/coverage"

# Add certificates for TLS requests
RUN apk --no-cache add ca-certificates

# Expose default port
EXPOSE 8080

# Copy the binary from the dev stage into a location that is in PATH
COPY --from=dev /go/bin/coverage /usr/local/bin/

# Set binary as the default command
CMD ["coverage"]
