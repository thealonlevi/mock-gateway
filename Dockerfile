######################## build stage ########################
FROM golang:1.22-alpine AS builder
WORKDIR /src

# Improve cache usage: fetch modules first
COPY go.mod go.sum ./
RUN go mod download

# Build-time argument: upstream host:port for sdk-server NLB
# e.g. "--build-arg UPSTREAM=sdk-server-nlb:9090"
ARG UPSTREAM=sdk-server-nlb:9090
ENV UPSTREAM=${UPSTREAM}

# Copy source and compile
COPY . .

# Inject upstream via -ldflags: -X main.upstream=<host:port>
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w -X main.upstream=${UPSTREAM}" \
    -o /bin/mock-gateway ./main.go

####################### runtime stage #######################
FROM gcr.io/distroless/static-debian11

# Minimal container with static binary
COPY --from=builder /bin/mock-gateway /mock-gateway

EXPOSE 8080
ENTRYPOINT ["/mock-gateway"]
