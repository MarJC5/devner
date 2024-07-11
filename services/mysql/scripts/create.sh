#!/bin/bash

echo "Please enter the database name:"
read dbname

echo "Please enter the database user name:"
read dbuser

echo "Please enter the database password:"
read dbpass

# Login to MySQL, you can replace root with another user with necessary privileges
mysql -u root -p'devner' <<EOF
CREATE DATABASE IF NOT EXISTS \`$dbname\` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
CREATE USER IF NOT EXISTS '$dbuser'@'%' IDENTIFIED BY '$dbpass';
GRANT ALL PRIVILEGES ON \`$dbname\`.* TO '$dbuser'@'%';
FLUSH PRIVILEGES;
EOF

# Check if the database and user were created
if [ $? -ne 0 ]; then
  echo "Error creating database and user."
  exit 1
fi

echo "Database and user created."
echo "Database: $dbname"
echo "Username: $dbuser"
