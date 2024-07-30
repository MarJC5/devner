#!/bin/bash

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# Function to display help
show_help() {
    echo "Usage: $0 <database_name> <database_user>"
    echo ""
    echo "This script checks for the existence of a MySQL database and user."
    echo ""
    echo "Arguments:"
    echo "  <database_name>     The name of the database to check."
    echo "  <database_user>     The name of the user to check."
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

# Check if the database exists
db_exists=$(mysql -u root -p'devner' -e "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '$dbname';" 2>/dev/null | grep -q $dbname && echo "yes" || echo "no")

# Check if the user exists
user_exists=$(mysql -u root -p'devner' -e "SELECT User FROM mysql.user WHERE User = '$dbuser';" 2>/dev/null | grep -q $dbuser && echo "yes" || echo "no")

# Output the results
if [ "$db_exists" = "yes" ]; then
    echo -e "${GREEN}Database '${dbname}' already exists.${NC}"
else
    echo -e "${RED}Database '${dbname}' does not exist.${NC}"
fi

if [ "$user_exists" = "yes" ]; then
    echo -e "${GREEN}User '${dbuser}' already exists.${NC}"
else
    echo -e "${RED}User '${dbuser}' does not exist.${NC}"
fi

# Exit with appropriate status
if [ "$db_exists" = "yes" ] && [ "$user_exists" = "yes" ]; then
    exit 0
else
    exit 1
fi

exit 0
