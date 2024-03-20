#!/bin/bash

# Color Codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

set -e

# FOLDER
SRCS_DIR=$(cd $(dirname $0); cd ..; pwd)
NGINX_SITES_AVAILABLE="${SRCS_DIR}/sites-enabled"
NGINX_SITES_CONFIG="${SRCS_DIR}/config/default.conf"
APPS_DIR="/home/dev/apps"

echo "Please enter the host name:"
read host

echo -e "${BLUE}Select a PHP version:${NC}"
options=("php_8_2" "php_8_1" "php_8" "php_7_4")
for i in ${!options[@]}; do
    echo -e "${YELLOW}$((i+1)))${NC} ${options[$i]}"
done

while true; do
    read -p "#? " num
    case $num in
        1|2|3|4)
            php_version=${options[$((num-1))]}
            echo -e "You selected ${GREEN}$php_version${NC}"
            break
            ;;
        *)
            echo -e "${RED}Invalid option, please select a valid number.${NC}"
            ;;
    esac
done


# Add the domain to the /etc/hosts file if it doesn't already exist
if ! grep -q "$host" /etc/hosts; then
    # Return line
    sudo -- sh -c -e "echo '' >> /etc/hosts"
    sudo -- sh -c -e "echo '127.0.0.1\t${host}' >> /etc/hosts" # IPv4
    sudo -- sh -c -e "echo '::1\t\t${host}' >> /etc/hosts" # IPv6
    sudo -- sh -c -e "echo '127.0.0.1\twww.${host}' >> /etc/hosts" # IPv4 www
    sudo -- sh -c -e "echo '::1\t\twww.${host}' >> /etc/hosts" # IPv6 www
    ech -e "${GREEN}$host${NC} added to /etc/hosts"
else
    echo -e "${YELLOW}$host${NC} host already existing in /etc/hosts."
fi

# Create Nginx virtual host file if it doesn't already exist
VHOST_FILE="${NGINX_SITES_AVAILABLE}/${host}.conf"

if [ ! -f "$VHOST_FILE" ]; then
    echo "Creating Nginx virtual host for $host..."
    
    sudo bash -c "cat > ${VHOST_FILE}" <<EOF
server {
    listen 80;
    root ${APPS_DIR}/${host};
    server_name ${host}.local www.${host}.local;

    index index.php index.html index.htm;

    location / {
        try_files \$uri \$uri/ /index.php?\$query_string;
    }

    include /etc/nginx/php.conf;
}
EOF

    echo -e "Virtual host for ${GREEN}$host${NC} created and enabled."
else
    echo -e "Virtual host file for ${YELLOW}$host${NC} already exists."
fi

if [ -f "$NGINX_SITES_CONFIG" ]; then
    escaped_host=$(echo "$host" | sed 's/\./\\./g')

    # Prepare the mapping line to insert
    mapping="    '$escaped_host' '$php_version';"
    

    # Check if the mapping for the host already exists
    if grep -q "$mapping" "$NGINX_SITES_CONFIG"; then
        echo -e "Mapping for ${YELLOW}$host${NC} already exists in the Nginx configuration."
    else
        echo "Updating Nginx default.conf file..."
        # Prepare the mapping line to insert, ensuring it ends with a newline
        if [[ "$OSTYPE" == "darwin"* ]]; then
            # macOS
            sudo perl -i -pe "s/}\$/\\n$mapping\\n}/" "$NGINX_SITES_CONFIG"
        else
            # GNU/Linux
            sudo sed -i -z "s/}\$/$mapping\n}/" "$NGINX_SITES_CONFIG"
        fi

        echo -e "Mapping for ${GREEN}$host${NC} added to the Nginx configuration."
    fi
fi
