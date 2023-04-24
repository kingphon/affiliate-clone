# Start from the latest golang base image
FROM golang:latest as builder

# Add Maintainer Info
LABEL maintainer="Selly <dev@selly.vn>"

# Set the Current Working Directory inside the container
RUN mkdir -p /go/src/git.selly.red/Selly-Server/affiliate
WORKDIR /go/src/git.selly.red/Selly-Server/affiliate

# Copy data to working dir
COPY . .

# Install dependencies
RUN go mod download

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo -o ./cmd/admin/main ./cmd/admin/main.go

######## Start a new stage from scratch #######
FROM alpine:latest

RUN apk --no-cache add tzdata zip ca-certificates

WORKDIR /go/src/git.selly.red/Selly-Server/affiliate

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /go/src/git.selly.red/Selly-Server/affiliate .

# Command to run the executable
CMD ["./cmd/admin/main"]