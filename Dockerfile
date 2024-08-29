# Stage 1: Build the frontend
FROM node:20.12.1-alpine as frontend-stage

# Set the working directory in the container
WORKDIR /app/frontend-bundle

# Copy frontend-bundle
COPY frontend-bundle ./

# Install dependencies
RUN npm install && npm run build

# Copy the rest of your app's source code
COPY . .

# Build the app
RUN npm run build

# Stage 1: Build Go application
FROM golang:1.22.3 as backend-stage

WORKDIR /app

COPY . .

ARG DB_CONNECTION=mysql

RUN echo "DB_CONNECTION is set to: $DB_CONNECTION"

# Conditional build step for SQLite or other DB connections
RUN go clean --modcache && go mod tidy && \
    if [ "$DB_CONNECTION" = "sqlite" ]; then \
        echo "Building with SQLite support"; \
        CGO_ENABLED=1 GOOS=linux go build -o main; \
    else \
        echo "Building without SQLite support"; \
        CGO_ENABLED=0 GOOS=linux go build -o main; \
    fi

# Stage 2: Build libwebp from source
FROM alpine:latest AS builder2

RUN apk --no-cache add \
    libpng-dev \
    libjpeg-turbo-dev \
    giflib-dev \
    tiff-dev \
    autoconf \
    automake \
    make \
    gcc \
    g++ \
    wget

RUN wget https://storage.googleapis.com/downloads.webmproject.org/releases/webp/libwebp-0.6.0.tar.gz && \
    tar -xvzf libwebp-0.6.0.tar.gz && \
    mv libwebp-0.6.0 libwebp && \
    rm libwebp-0.6.0.tar.gz && \
    cd /libwebp && \
    ./configure && \
    make && \
    make install && \
    rm -rf libwebp

# Stage 3: Final stage with minimal image
FROM alpine:latest

COPY --from=builder2 /usr/local/bin /usr/local/bin
COPY --from=builder2 /usr/local/include /usr/local/include
COPY --from=builder2 /usr/local/lib /usr/local/lib

RUN apk --no-cache add libpng libjpeg-turbo giflib tiff && \
    rm -rf /usr/local/share /usr/local/libexec

WORKDIR /root/

COPY --from=backend-stage ./app/main .
COPY --from=frontend-stage ./app/static ./static

EXPOSE 8000

CMD ["./main"]
