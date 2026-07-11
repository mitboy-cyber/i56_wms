#!/bin/bash
# Run PostgreSQL migrations for i56_dev
set -e

DB="i56_dev"
MIGDIR="/home/ubuntu/i56/migrations"

echo "=== Running I56 PostgreSQL Migrations ==="

for f in $(ls "$MIGDIR"/*.up.sql | sort); do
    echo "Running: $(basename "$f")"
    psql -d "$DB" -f "$f" -q
done

echo "=== All migrations complete ==="

# Verify tables
echo ""
echo "=== Tables created ==="
psql -d "$DB" -c "\dt"

echo ""
echo "=== Row counts ==="
for t in tenants parcels orders warehouses clients carriers routes roles users; do
    count=$(psql -d "$DB" -t -c "SELECT COUNT(*) FROM $t" 2>/dev/null || echo "N/A")
    echo "  $t: $count rows"
done
