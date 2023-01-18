#!/usr/bin/env bash

name=code-social-mysql
image=mysql
username=code-social
password=password
dbname=code_social

[[ $(docker ps -f "name=$name" --format '{{.Names}}') == $name ]] || docker run -d \
    --name "$name" \
    -p 3306:3306 \
    -e MYSQL_ROOT_PASSWORD=root \
    -e MYSQL_USER=$username \
    -e MYSQL_PASSWORD=$password \
    -e MYSQL_DATABASE=$dbname \
    $image
