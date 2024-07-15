#!/bin/bash

echo "Please enter the database name:"
read dbname

echo "Please enter the database user name:"
read dbuser

echo "Please enter the database password:"
read -s dbpass

# Check if dblink extension is available
psql -U devner -d devner -c "CREATE EXTENSION IF NOT EXISTS dblink;"

# Login to PostgreSQL as the postgres user and create the database if it doesn't exist
psql -U devner -d devner -tc "SELECT 1 FROM pg_database WHERE datname = '$dbname'" | grep -q 1 || psql -U devner -c "CREATE DATABASE $dbname WITH ENCODING 'UTF8';"

# Create the user if it doesn't exist
psql -U devner -d devner <<EOF
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
