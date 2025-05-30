#!/bin/bash

API_PORT=${API_PORT:-8000}
BASE_URL="http://localhost:$API_PORT"

echo "Getting all quotes from $BASE_URL..."

curl $BASE_URL/quotes
echo ""

echo "Done."