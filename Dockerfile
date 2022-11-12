# Vendor stage
FROM golang:1.19.2 as dep
WORKDIR /build
COPY go.mod go.sum ./
COPY . .
RUN go mod vendor

# Build binary stage
FROM golang:1.19.2 as build
WORKDIR /build
COPY --from=dep /build .
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o server -tags nethttpomithttp2 ./cmd

# Minimal image
FROM alpine:latest
WORKDIR /app
COPY internal/pkg/config/file internal/pkg/config/file
COPY migrations/mysql migrations/mysql
COPY --from=build /build/server server
RUN apk update
RUN apk upgrade
RUN apk add ca-certificates
RUN apk --no-cache add tzdata
CMD ["./server"]
