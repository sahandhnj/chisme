#!/bin/bash

# Base URL
BASE_URL="http://localhost:4004"

# 1. Get All Servers
echo "Fetching all servers..."
curl -X GET "$BASE_URL/mock/servers"
echo -e "\n"

# 2. Get Specific Server by ID
SERVER_ID=1
echo "Fetching server with ID $SERVER_ID..."
curl -X GET "$BASE_URL/mock/server/$SERVER_ID"
echo -e "\n"

# 3. Get Specific Application by Server ID and Application ID
APP_ID=2
echo "Fetching application with ID $APP_ID on server with ID $SERVER_ID..."
curl -X GET "$BASE_URL/mock/server/$SERVER_ID/application/$APP_ID"
echo -e "\n"