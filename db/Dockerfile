# Build Redis from source to avoid Docker Hub dependency
FROM ubuntu:20.04

# Avoid interactive prompts during build
ENV DEBIAN_FRONTEND=noninteractive

# Install dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    wget \
    && rm -rf /var/lib/apt/lists/*

# Download and compile Redis
RUN wget http://download.redis.io/redis-stable.tar.gz \
    && tar xzf redis-stable.tar.gz \
    && cd redis-stable \
    && make \
    && make install \
    && cd .. \
    && rm -rf redis-stable redis-stable.tar.gz

# Create redis user and data directory
RUN useradd -r -s /bin/false redis \
    && mkdir -p /data \
    && chown redis:redis /data

# Switch to redis user
USER redis

# Set working directory
WORKDIR /data

# Expose Redis port
EXPOSE 6379

# Start Redis server
CMD ["redis-server", "--bind", "0.0.0.0", "--port", "6379", "--dir", "/data", "--appendonly", "yes"]