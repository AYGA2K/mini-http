#!/bin/bash

BASE_URL="http://localhost:8080"

echo "Testing GET /"
curl -i "$BASE_URL/"

echo -e "\n\nTesting GET /user/ayoub"
curl -i "$BASE_URL/user/ayoub"

echo -e "\n\nTesting GET /search?q=golang"
curl -i "$BASE_URL/search?q=golang"

echo -e "\n\nTesting POST /submit"
curl -i -X POST "$BASE_URL/submit" -d "Hello body"

echo -e "\n\nTesting GET /submit (should return 405)"
curl -i "$BASE_URL/submit"
