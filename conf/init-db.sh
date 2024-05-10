#!/bin/bash

set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE USER replication REPLICATION LOGIN PASSWORD 'replication';
    GRANT REPLICATION CLIENT ON *.* TO replication;
EOSQL
