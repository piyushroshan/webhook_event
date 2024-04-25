FROM golang:1.21.0

WORKDIR /app

COPY go.mod ./
COPY src/ ./
RUN ls -la

# Build
RUN go build -o main .

EXPOSE 9999

# Run
CMD ["/app/main"]