#!/usr/bin/env bash

name=code-things-pg
image=postgres
username=code-things
password=password
dbname=code_things
port=5432

[[ $(docker ps -f "name=$name" --format '{{.Names}}') == $name ]] || docker run -d \
    --name "$name" \
    -p $port:$port \
    -e POSTGRES_USER=$username \
    -e POSTGRES_PASSWORD=$password \
    -e POSTGRES_DB=$dbname \
    -v "$(pwd)/scripts/db/data:/var/lib/postgresql/data" \
    -v "$(pwd)/scripts/db/initdb.d:/docker-entrypoint-initdb.d" \
    $image
