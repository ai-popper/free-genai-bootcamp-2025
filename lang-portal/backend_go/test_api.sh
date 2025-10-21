#!/bin/bash

# Test script for Language Learning Portal API

BASE_URL="http://localhost:8080"

echo "=== Language Learning Portal API Tests ==="
echo ""

echo "1. Health Check"
curl -s $BASE_URL/health | jq .
echo -e "\n"

echo "2. Get All Groups"
curl -s $BASE_URL/api/groups | jq '.items[] | {id, name, word_count}'
echo -e "\n"

echo "3. Get Words from Group 1 (Basic Greetings)"
curl -s "$BASE_URL/api/groups/1/words" | jq '.items[0:3] | .[] | {japanese, romaji, english}'
echo -e "\n"

echo "4. Get Study Activities"
curl -s $BASE_URL/api/study_activities/1 | jq .
echo -e "\n"

echo "5. Create Study Session"
SESSION_RESPONSE=$(curl -s -X POST $BASE_URL/api/study_activities \
  -H "Content-Type: application/json" \
  -d '{"group_id": 2, "study_activity_id": 1}')
echo $SESSION_RESPONSE | jq .
SESSION_ID=$(echo $SESSION_RESPONSE | jq -r '.id')
echo -e "\n"

echo "6. Record Word Reviews"
curl -s -X POST "$BASE_URL/api/study_sessions/$SESSION_ID/words/11/review" \
  -H "Content-Type: application/json" \
  -d '{"correct": true}' | jq .
echo ""

curl -s -X POST "$BASE_URL/api/study_sessions/$SESSION_ID/words/12/review" \
  -H "Content-Type: application/json" \
  -d '{"correct": false}' | jq .
echo -e "\n"

echo "7. Get Dashboard Stats"
curl -s $BASE_URL/api/dashboard/quick_stats | jq .
echo -e "\n"

echo "8. Get Study Progress"
curl -s $BASE_URL/api/dashboard/study_progress | jq .
echo -e "\n"

echo "9. Get Last Study Session"
curl -s $BASE_URL/api/dashboard/last_study_session | jq .
echo -e "\n"

echo "10. Get All Study Sessions"
curl -s $BASE_URL/api/study_sessions | jq '.items[] | {id, activity_name, group_name, review_items_count}'
echo -e "\n"

echo "=== Tests Complete ==="
