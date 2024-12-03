FROM golang:1.22

# Set destination for COPY
WORKDIR /app

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/reference/dockerfile/#copy
COPY . .
RUN go mod tidy
RUN go mod download

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /concerts-api

EXPOSE 8080

# Run
CMD ["/concerts-api"]
