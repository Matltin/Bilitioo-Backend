#!/bin/sh

set -e

echo "run db migration"
/app/migrate -path /app/migration -datanase "$DB_SOURCE" -verbose up

echo "strat the app"
exec "$@"