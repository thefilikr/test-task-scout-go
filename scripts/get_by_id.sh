#!/bin/bash

# ./scripts/get_by_id.sh "quote-id"

API_PORT=${API_PORT:-8000}
BASE_URL="http://localhost:$API_PORT"

QUOTE_ID="$1"

if [ -z "$QUOTE_ID" ]; then
  echo "Usage: $0 \"quote-id\""
  echo "You can get quote IDs by running ./scripts/get_all_quotes.sh"
  exit 1
fi

echo "Getting quote with ID: $QUOTE_ID from $BASE_URL..."

curl $BASE_URL/quotes/$QUOTE_ID
echo ""

echo "Done." 