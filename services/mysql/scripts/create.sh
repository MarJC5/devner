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
    echo "Usage: $0 <database_name> <database_user> <database_password>"
    echo ""
    echo "This script creates a MySQL database and user with the specified credentials."
    echo ""
    echo "Arguments:"
    echo "  <database_name>     The name of the database to create."
    echo "  <database_user>     The name of the user to create."
    echo "  <database_password> The password for the new user."
}

# Check if exactly 3 arguments are passed
if [ $# -ne 3 ]; then
    show_help
    exit 1
fi

dbname=$1
dbuser=$2
dbpass=$3

# Format the database name and user
dbname=$(echo $dbname | sed 's/[^a-zA-Z0-9]//g')
dbuser=$(echo $dbuser | sed 's/[^a-zA-Z0-9]//g')

# Check if the database already exists
if [ $(mysql -u root -p'devner' -e "SELECT SCHEMA_NAME FROM INFORMATION_SCHEMA.SCHEMATA WHERE SCHEMA_NAME = '$dbname';" | wc -l) -ne 1 ]; then
    echo -e "${RED}Database already exists.${NC}"
    exit 1
fi

# Login to MySQL and execute commands
mysql -u root -p'devner' <<EOF
CREATE DATABASE IF NOT EXISTS \`$dbname\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '$dbuser'@'%' IDENTIFIED BY '$dbpass';
GRANT ALL PRIVILEGES ON \`$dbname\`.* TO '$dbuser'@'%';
FLUSH PRIVILEGES;
EOF

# Check if the database and user were created
if [ $? -ne 0 ]; then
  echo -e "${RED}Error creating database and user.${NC}"
  exit 1
fi

echo -e "${GREEN}Database and user created.${NC}"
echo "Database: $dbname"
echo "Username: $dbuser"

