#!/bin/sh
set -eu

cqlsh cassandra -f /init.cql

# CREATE TABLE IF NOT EXISTS does not update tables from older deployments.
columns=$(cqlsh cassandra -e "SELECT column_name FROM system_schema.columns WHERE keyspace_name = 'todolist' AND table_name = 'tasks' AND column_name = 'title';")
if ! printf '%s\n' "$columns" | grep -Eq '^[[:space:]]*title[[:space:]]*$'; then
    echo "Adding the task title column to the existing schema"
    cqlsh cassandra -e "ALTER TABLE todolist.tasks ADD title text;"
fi
