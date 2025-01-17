# Build stage
FROM golang:1.23.5-alpine as build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o sync ./cmd/main

# Final stage
FROM --platform=$TARGETPLATFORM alpine:3.21
WORKDIR /app

COPY --from=build /app/sync .

RUN apk --no-cache add ca-certificates tzdata git openssh-client

ENTRYPOINT ["/app/sync"]
