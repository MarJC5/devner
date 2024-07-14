#!/bin/bash

# Function to check if a database exists
function db_exists() {
  local db=$1
  psql -U postgres -tAc "SELECT 1 FROM pg_database WHERE datname='$db'" | grep -q 1
}

# Function to check if a user exists
function user_exists() {
  local user=$1
  psql -U postgres -tAc "SELECT 1 FROM pg_roles WHERE rolname='$user'" | grep -q 1
}

echo "Please enter the database name:"
read dbname

echo "Please enter the database user name:"
read dbuser

echo "Please enter the database password:"
read -s dbpass

# Check if database exists
if db_exists $dbname; then
  echo "Database '$dbname' already exists."
else
  # Create the database
  createdb -U postgres -E UTF8 $dbname
  if [ $? -ne 0 ]; then
    echo "Error creating database '$dbname'."
    exit 1
  fi
  echo "Database '$dbname' created."
fi

# Check if user exists
if user_exists $dbuser; then
  echo "User '$dbuser' already exists."
else
  # Create the user
  psql -U postgres -c "CREATE USER $dbuser WITH PASSWORD '$dbpass';"
  if [ $? -ne 0 ]; then
    echo "Error creating user '$dbuser'."
    exit 1
  fi
  echo "User '$dbuser' created."
fi

# Grant all privileges on the database to the user
psql -U postgres -c "GRANT ALL PRIVILEGES ON DATABASE $dbname TO $dbuser;"
if [ $? -ne 0 ]; then
  echo "Error granting privileges to user '$dbuser' on database '$dbname'."
  exit 1
fi

echo "All privileges on database '$dbname' granted to user '$dbuser'."

echo "Database and user setup completed."
echo "Database: $dbname"
echo "Username: $dbuser"
