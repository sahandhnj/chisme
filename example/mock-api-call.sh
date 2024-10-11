#!/bin/bash

# Base URL
BASE_URL="localhost:4004"

# 1. Get All Servers
echo "Fetching all servers..."
curl -X GET "$BASE_URL/mock/servers"
echo -e "\n"

# 2. Get All Applications
echo "Fetching all applications..."
curl -X GET "$BASE_URL/mock/applications"
echo -e "\n"

# 3. Get All Resources
echo "Fetching all resources..."
curl -X GET "$BASE_URL/mock/resources"
echo -e "\n"

# 4. Get Specific Server by ID
SERVER_ID=1
echo "Fetching server with ID $SERVER_ID..."
curl -X GET "$BASE_URL/mock/server/$SERVER_ID"
echo -e "\n"

# 5. Get Specific Application by Server ID and Application ID
APP_ID=2
echo "Fetching application with ID $APP_ID on server with ID $SERVER_ID..."
curl -X GET "$BASE_URL/mock/server/$SERVER_ID/application/$APP_ID"
echo -e "\n"

# 6. Get Specific Resource by Server ID and Resource ID
RESOURCE_ID=6
echo "Fetching application with ID $RESOURCE_IDon server with ID $SERVER_ID..."
curl -X GET "$BASE_URL/mock/server/$SERVER_ID/resource/$RESOURCE_ID"
echo -e "\n"

