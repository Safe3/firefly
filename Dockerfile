FROM alpine:latest

LABEL org.opencontainers.image.title="firefly" \
      org.opencontainers.image.version="v4.3" \
      org.opencontainers.image.description="Firefly WireGuard Server" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.source="https://github.com/Safe3/firefly"

COPY firefly-linux-amd64 /firefly/firefly

# Install Linux packages
RUN apk add --no-cache --purge --clean-protected dumb-init iptables tzdata && rm -rf /var/cache/apk/*

EXPOSE 50120/udp
EXPOSE 50121/tcp

WORKDIR /firefly
CMD ["/usr/bin/dumb-init", "./firefly"]
