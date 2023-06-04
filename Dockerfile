# Use a Golang base image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

ENV API_PORT=8080
ENV ETH_NODE_URL="https://eth-goerli.g.alchemy.com/v2"
ENV DB_CONNECTION_URL="postgres://postgres:postgres@postgres:5432/postgres?sslmode=disable"
ENV API_KEY="jEvj-KdZ92ZUmX01Jpegiu52fpgEpE8_"
ENV JWT_SECRET="secret"

# Copy the source code to the working directory
COPY . .

# Build the Go application
RUN go build -o limeapi .

# Expose the port that your server listens on
EXPOSE 8080

# Set the entry point to run the Lime API server
CMD ["./limeapi"]

#docker run -p 8080:8080 -e API_PORT=8080 -e ETH_NODE_URL="https://eth-goerli.g.alchemy.com/v2" -e DB_CONNECTION_URL="postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" limeapi10