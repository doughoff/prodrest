# Start from a Debian based image with Go installed
FROM golang:1.20.4-buster

# Set the Current Working Directory inside the container
WORKDIR /app

# We add the go mod and go sum files before the rest of the code
# to leverage Docker cache
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy everything from the current directory to the PWD(Present Working Directory) inside the container
COPY . .

# Install the package
RUN go build -o bin/prodapi .

# This container exposes port 8080 to the outside world
EXPOSE 8080

# Run the binary program produced by `go install`
CMD ["./bin/prodapi"]
