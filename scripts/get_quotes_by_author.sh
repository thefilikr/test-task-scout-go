#!/bin/bash
# ./scripts/get_quotes_by_author.sh "Author Name"

API_PORT=${API_PORT:-8000}
BASE_URL="http://localhost:$API_PORT"

AUTHOR="$1"

if [ -z "$AUTHOR" ]; then
  echo "Usage: $0 \"Author Name\""
  exit 1
fi

echo "Getting quotes by author: \"$AUTHOR\" from $BASE_URL..."

curl "$BASE_URL/quotes?author=$(echo "$AUTHOR" | sed 's/ /%20/g')"
echo ""

echo "Done."