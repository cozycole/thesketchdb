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
    psql -d $1 -U "$USER" -c "
        CREATE EXTENSION IF NOT EXISTS citext;
        CREATE TYPE cast_role AS ENUM ('cast', 'guest', 'host', '');
        CREATE TYPE character_type AS ENUM ('original', 'impression', 'fictional_impression', 'generic');
    " 
    if [ $? -ne 0 ]; then
        exit 1
    fi
fi
