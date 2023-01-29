# Build api-gateway micro service image
FROM golang:alpine3.17 AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create appuser.
RUN adduser -D -g '' appuser

WORKDIR $GOPATH/src/github.com/iamvasanth07/showcase/api-gateway

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/api-gateway .

# Final stage: the running container.
FROM scratch
# Import the api-gateway and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable.
COPY --from=builder /go/bin/api-gateway /go/bin/api-gateway

# Use an unprivileged api-gateway.
USER appuser

# Run the app.
ENTRYPOINT ["/go/bin/api-gateway"]



