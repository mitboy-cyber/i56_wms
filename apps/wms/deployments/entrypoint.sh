#!/bin/sh
set -e

echo "I56 WMS — Starting"

# Start Go backend
/opt/i56/i56-server &
GO_PID=$!

# Start nginx
nginx -g "daemon off;" &
NGINX_PID=$!

echo "I56 WMS — Running (Go: $GO_PID, Nginx: $NGINX_PID)"

# Wait for either to exit
wait -n
echo "Process exited, shutting down..."
kill $GO_PID $NGINX_PID 2>/dev/null || true
