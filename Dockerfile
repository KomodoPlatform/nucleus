# Compilation
FROM docker.io/library/golang:1.21-alpine AS nucleus-builder
WORKDIR /src/src/nucleus
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
ENV PACKAGES curl make git libc-dev bash gcc linux-headers eudev-dev python3
RUN apk add --no-cache $PACKAGES
RUN CGO_ENABLED=0 make install

# Runtime
FROM docker.io/chainguard/static:latest
COPY --from=nucleus-builder /go/bin/nucleusd /usr/local/bin/
USER 0

ENTRYPOINT ["nucleusd", "start"]
