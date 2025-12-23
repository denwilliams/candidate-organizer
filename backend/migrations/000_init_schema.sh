#!/bin/bash
set -e

# Get schema name from environment variable, default to 'public'
SCHEMA_NAME="${POSTGRES_SCHEMA:-public}"

echo "Initializing database schema: $SCHEMA_NAME"

# Create schema if it doesn't exist (skip if using public schema)
if [ "$SCHEMA_NAME" != "public" ]; then
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
        CREATE SCHEMA IF NOT EXISTS $SCHEMA_NAME;
        GRANT ALL ON SCHEMA $SCHEMA_NAME TO $POSTGRES_USER;
EOSQL
    echo "Schema $SCHEMA_NAME created successfully"
else
    echo "Using default public schema"
fi
