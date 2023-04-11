<<<<<<< HEAD
<<<<<<< HEAD
<<<<<<< HEAD
FROM golang:1.18.8-buster
=======
FROM golang:1.18.9-buster
>>>>>>> bdc8cc7db62ce10f1b96b5df06804636a421f2a1
=======
FROM golang:1.18.9-buster
>>>>>>> 96309aa37d5874e97035fe95a3c08dab0be7fdfe
=======
FROM golang:1-buster
>>>>>>> 1a597a954d6c46633ff58aa9aca8e76f37e8c68a

# Set the Current Working Directory inside the container
WORKDIR /goapp

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

EXPOSE 8080

# Build the Go app
RUN go build -o gh-mgnt .

# Command to run the executable
CMD ./gh-mgnt