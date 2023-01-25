# Build user micro service 
FROM golang:1.12.5-alpine3.9 as builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git

# Create appuser.
RUN adduser -D -g '' appuser

WORKDIR $GOPATH/src/github.com/iamvasanth07/showcase/user
COPY . .

# Fetch dependencies.
# Using go get.
RUN go mod download

# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/app .

# Final stage: the running container.
FROM scratch
# Import the user and group files from the builder.
COPY --from=builder /etc/passwd /etc/passwd
# Copy our static executable.
COPY --from=builder /go/bin/app /go/bin/app

# Use an unprivileged user.
USER appuser

# Expose port 8080 to the outside world
EXPOSE 8080

# Run the app.
ENTRYPOINT ["/go/bin/app"]



