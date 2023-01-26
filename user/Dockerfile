# Build user micro service image
FROM golang:alpine3.17 AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create appuser.
RUN adduser -D -g '' appuser

WORKDIR $GOPATH/src/github.com/iamvasanth07/showcase/user

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/user .

# Final stage: the running container.
FROM scratch
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable.
COPY --from=builder /go/bin/user /go/bin/user

# Use an unprivileged user.
USER appuser

# Run the app.
CMD ["/go/bin/./user"]



