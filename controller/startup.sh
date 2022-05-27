#!/bin/sh

# Startup Redis
echo "Starting redis..."
redis-server --version
redis-server /etc/redis/redis.conf --daemonize yes

# Startup controller
echo "Starting controller..."
REDIS_HOST=0.0.0.0 REDIS_PORT=6379 GROUP_TASKS_FILE=spec_context_map.json controller