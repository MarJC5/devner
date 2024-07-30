#!/bin/bash

# Colors
RED='\033[0;31m'
BRIGHT_RED='\033[0;91m'
CYAN='\033[0;36m'
BRIGHT_CYAN='\033[0;96m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to display help
show_help() {
    echo "Usage: $0 <database_name> <database_user>"
    echo ""
    echo "This script removes a PostgreSQL database and user with the specified credentials."
    echo ""
    echo "Arguments:"
    echo "  <database_name>     The name of the database to remove."
    echo "  <database_user>     The name of the user to remove."
}

# Check if exactly 2 arguments are passed
if [ $# -ne 2 ]; then
    show_help
    exit 1
fi

dbname=$1
dbuser=$2

# Format the database name and user
dbname=$(echo $dbname | sed 's/[^a-zA-Z0-9]//g')
dbuser=$(echo $dbuser | sed 's/[^a-zA-Z0-9]//g')

# Revoke connections from the database
psql -U devner -d devner -c "REVOKE CONNECT ON DATABASE $dbname FROM public;"

# Terminate all existing connections to the database
psql -U devner -d devner -c "SELECT pg_terminate_backend(pg_stat_activity.pid) FROM pg_stat_activity WHERE pg_stat_activity.datname = '$dbname' AND pid <> pg_backend_pid();"

# Drop the database if it exists
psql -U devner -d devner -c "DROP DATABASE IF EXISTS $dbname;"

# Drop the user if it exists
psql -U devner -d devner -c "DROP USER IF EXISTS $dbuser;"

# Check if the database and user were removed successfully
if [ $? -ne 0 ]; then
  echo "Error removing database and user."
  exit 1
fi

echo -e "${GREEN}Database and user removed.${NC}"

exit 0