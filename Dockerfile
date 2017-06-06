FROM golang:latest
# Expose default port
EXPOSE 3000

# gin is for dev monitoring
# RUN go-wrapper download github.com/codegangsta/gin
# RUN go-wrapper install github.com/codegangsta/gin

# Copy the local package files to the container’s workspace.
ADD . /go/src/github.com/archivers-space/coverage
# WORKDIR /go/src/github.com/archivers-space/coverage
# CMD ["gin", "-i"]

# Install api binary globally within container 
RUN go install github.com/archivers-space/coverage
# Set binary as entrypoint
ENTRYPOINT /go/bin/coverage



