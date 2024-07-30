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
    echo "This script creates a PostgreSQL database and user with the specified credentials."
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
if [ $(psql -U devner -d devner -tc "SELECT 1 FROM pg_database WHERE datname = '$dbname'" | grep -q 1) ]; then
    echo -e "${RED}Database already exists.${NC}"
    exit 1
fi

# Check if dblink extension is available
psql -U devner -d devner -c "CREATE EXTENSION IF NOT EXISTS dblink;"

# Login to PostgreSQL as the devner user and create the database if it doesn't exist
psql -U devner -d devner -tc "SELECT 1 FROM pg_database WHERE datname = '$dbname'" | grep -q 1 || psql -U devner -c "CREATE DATABASE $dbname WITH ENCODING 'UTF8';"

# Create the user if it doesn't exist and grant privileges
psql -U devner -d devner <<EOF
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '$dbuser') THEN
      EXECUTE 'CREATE USER $dbuser WITH PASSWORD ''$dbpass''';
   END IF;
END
\$\$;

GRANT ALL PRIVILEGES ON DATABASE $dbname TO $dbuser;
EOF

# Check if the database and user were created successfully
if [ $? -ne 0 ]; then
  echo -e "${RED}Error creating database and user.${NC}"
  exit 1
fi

echo -e "${GREEN}Database and user created.${NC}"
echo "Database: $dbname"
echo "Username: $dbuser"
