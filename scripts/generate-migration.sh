#!/usr/bin/env bash

if [ -z "$1" ]; then
    echo "Error: Please provide a migration name as an argument."
    exit 1
fi

migrations_dir="../pkg/state/migrations"

if [ ! -d "$migrations_dir" ]; then
    mkdir -p "$migrations_dir"
    if [ $? -ne 0 ]; then
        echo "Error: Could not create the migrations directory '$migrations_dir'."
        exit 1
    fi
fi

if [ ! -w "$migrations_dir" ]; then
    echo "Error: The migrations directory '$migrations_dir' is not writable."
    exit 1
fi

timestamp=$(date +"%Y%m%d%H%M%S")

migration_file="${migrations_dir}/${timestamp}_${1}.sql"

if [ -f "$migration_file" ]; then
    echo "Error: '$migration_file' already exists."
    exit 1
fi

touch "$migration_file" || {
    echo "Error: Could not create '$migration_file'. Check permissions."
    exit 1
}

if [ $? -ne 0 ]; then
    echo "Error: Failed to create migration file '$migration_file'."
    exit 1
fi

echo "Migration file '$migration_file' created successfully."
