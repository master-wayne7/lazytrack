FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1001 -S lazytrack && \
    adduser -u 1001 -S lazytrack -G lazytrack

# Set working directory
WORKDIR /app

# Copy binary
COPY lazytrack /app/lazytrack

# Make binary executable
RUN chmod +x /app/lazytrack

# Switch to non-root user
USER lazytrack

# Expose port (if needed in future)
# EXPOSE 8080

# Set entrypoint
ENTRYPOINT ["/app/lazytrack"]

# Default command
CMD ["--help"] 