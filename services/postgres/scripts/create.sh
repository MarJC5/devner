#!/bin/bash

echo "Please enter the database name:"
read dbname

echo "Please enter the database user name:"
read dbuser

echo "Please enter the database password:"
read -s dbpass

# Check if dblink extension is available
psql -U postgres -d postgres -c "CREATE EXTENSION IF NOT EXISTS dblink;"

# Login to PostgreSQL as the postgres user and execute the commands
psql -U postgres -d postgres <<EOF
-- Create the database if it doesn't exist
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = '$dbname') THEN
      EXECUTE 'CREATE DATABASE $dbname WITH ENCODING ''UTF8''';
   END IF;
END
\$\$;

-- Create the user if it doesn't exist
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '$dbuser') THEN
      EXECUTE 'CREATE USER $dbuser WITH PASSWORD ''$dbpass''';
   END IF;
END
\$\$;

-- Grant all privileges on the database to the user
GRANT ALL PRIVILEGES ON DATABASE $dbname TO $dbuser;
EOF

# Check if the database and user were created successfully
if [ $? -ne 0 ]; then
  echo "Error creating database and user."
  exit 1
fi

echo "Database and user created."
echo "Database: $dbname"
echo "Username: $dbuser"
