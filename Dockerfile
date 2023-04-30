FROM golang:1.19-alpine AS Builder
WORKDIR /app

# Download Depedency
COPY go.* ./
RUN go mod download

# Copy Source Code to Container
COPY . ./

# Build binary
RUN go build -v -o passenger-service

FROM alpine:3.9

# Copy binary from build stage
COPY --from=Builder /app/passenger-service /passenger-service

# Run app
CMD ["./passenger-service"]
