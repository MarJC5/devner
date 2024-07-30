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
    echo "This script removes a MySQL database and user with the specified credentials."
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

# Login to MySQL and execute commands to drop the database and user
mysql -u root -p'devner' 2>/dev/null <<EOF
DROP DATABASE IF EXISTS \`$dbname\`;
DROP USER IF EXISTS '$dbuser'@'%';
FLUSH PRIVILEGES;
EOF

# Check if the database and user were removed successfully
if [ $? -ne 0 ]; then
  echo -e "${RED}Error removing database and user.${NC}"
  exit 1
fi

echo -e "${GREEN}Database and user removed.${NC}"

exit 0