#!/bin/bash

API_PORT=${API_PORT:-8000}
BASE_URL="http://localhost:$API_PORT"

echo "Creating example quotes at $BASE_URL..."

curl -s -X POST \
  $BASE_URL/quotes \
  -H "Content-Type: application/json" \
  -d '{
    "text": "За свою улетность денег не беру, а за красоту тем более...",
    "author": "Панда По"
  }'
echo "" 

curl -s -X POST \
  $BASE_URL/quotes \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Счастье для всех, даром, и пусть никто не уйдет обиженный!",
    "author": "Редрик (Пикник на обочине)"
  }'
echo ""

curl -s -X POST \
  $BASE_URL/quotes \
  -H "Content-Type: application/json" \
  -d '{
    "text": "Вы всё твердите про белые и чёрные полосы, а я считаю, что даже все оттенки серого не смогут описать всю цветную красоту нашего мира!",
    "author": "The FILIkR"
  }'
echo ""

echo "Done creating quotes."