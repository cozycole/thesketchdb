#!/bin/bash 

if [ -z "$1" ]; then
    echo "Error: No database name provided"
    exit 1
else
    psql -d postgres -U "$USER" -c "CREATE DATABASE $1
        WITH
        OWNER = $USER
        ENCODING = 'UTF8'
        CONNECTION LIMIT = -1
        IS_TEMPLATE = False;
        "
    
    if [ $? -ne 0 ]; then
        exit 1
    fi

    psql -d $1 -U "$USER" -c "
        CREATE TYPE rating AS ENUM ('PG', 'PG-13', 'R');
    "
fi