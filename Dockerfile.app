# Build stage
FROM golang:1.22 AS builder

WORKDIR /app

# Copy go module files and download dependencies for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the binary with CGO disabled to create a static binary
RUN CGO_ENABLED=0 go build -v -o main

# Final stage
FROM alpine:3.15.11

# hadolint ignore=DL3018
RUN apk --no-cache add postgresql-client openjdk11-jre bash wget

# hadolint ignore=DL3047
RUN wget --quiet https://archive.apache.org/dist/kafka/3.4.0/kafka_2.13-3.4.0.tgz -O /tmp/kafka.tgz \
    && tar -xzf /tmp/kafka.tgz -C /opt \
    && rm /tmp/kafka.tgz

# Configure environment variables
ENV PATH="/opt/kafka_2.13-3.4.0/bin:${PATH}"
# Uncomment if Java is needed in the runtime environment
# ENV JAVA_HOME=/usr/lib/jvm/java-11-openjdk
ENV CLASSPATH="/opt/kafka_2.13-3.4.0/libs/*"

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy db_check.sh script and make it executable
COPY ./scripts/db_check.sh /usr/local/bin/db_check.sh
RUN chmod +x /usr/local/bin/db_check.sh

CMD ["./main"]
