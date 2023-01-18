#!/usr/bin/env bash

scheme="http"
host="127.0.0.1:3000"

echo "Creating the profile"
curl -v -X POST $scheme://$host/api/profile \
    -H 'Content-Type: application/json' \
    -d @scripts/create-profile.json
echo "------------------------------------------------------------"; echo

echo "Fetching the profile"
curl -v -X GET $scheme://$host/api/profile/justin
echo "------------------------------------------------------------"; echo

echo "Updating the avatar"
curl -v -X PUT $scheme://$host/api/profile/justin \
    -H 'Content-Type: application/json' \
    -d @scripts/update-profile.json
echo "------------------------------------------------------------"; echo

echo "Deleting profile"
curl -v -X DELETE $scheme://$host/api/profile/justin
echo "------------------------------------------------------------"; echo

echo "Fetching after delete should be 404"
curl -v -X GET $scheme://$host/api/profile/justin
echo "------------------------------------------------------------"; echo
