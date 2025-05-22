#!/bin/bash

# Path to the Caddyfile
CADDYFILE_PATH="${SRCS_DIR}/services/frankenphp/Caddyfile"

# Determine the hosts file location based on OS
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    HOSTS_FILE="/c/Windows/System32/drivers/etc/hosts"
else
    HOSTS_FILE="/etc/hosts"
fi

# Check if the script is running with sufficient privileges
check_privileges() {
    if [[ "$EUID" -ne 0 ]]; then
        echo -e "${RED}This script must be run as root or with administrative privileges.${NC}"
        exit 1
    fi
}

# Add a new entry to the hosts file
add_host_to_hosts_file() {
    local host=$1

    # Check if the host already exists
    if grep -q "$host" "$HOSTS_FILE"; then
        echo -e "${RED}Host $host already exists in the hosts file.${NC}"
        return 1
    fi

    # Append the new host entry
    echo "127.0.0.1 $host" >> "$HOSTS_FILE"
    echo "::1 $host" >> "$HOSTS_FILE"
    echo -e "${GREEN}Host $host added to the hosts file successfully.${NC}"
    return 0
}

# Remove an entry from the hosts file
remove_host_from_hosts_file() {
    local host=$1

    # Check if the host exists in the hosts file
    if ! grep -q "$host" "$HOSTS_FILE"; then
        echo -e "${RED}Host $host not found in the hosts file.${NC}"
        return 1
    fi

    # Use sed to remove the line
    if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
        sed -i "/$host/d" "$HOSTS_FILE"
    else
        sed -i "/$host/d" "$HOSTS_FILE"
    fi

    echo -e "${GREEN}Host $host removed from the hosts file successfully.${NC}"
    return 0
}

show_caddyfile_help() {
    echo -e "Usage: ${BOLD}$0 {add-host|remove-host|list-hosts} <host> [<root>${NC}]\n"
    echo -e "Commands:"
    echo -e "  ${YELLOW}add-host <host> <root>${NC}    Add a new host configuration."
    echo -e "  ${YELLOW}remove-host <host>${NC}        Remove an existing host configuration."
    echo -e "  ${YELLOW}list-hosts${NC}                List all hosts in the Caddyfile."
    echo -e "\nArguments:"
    echo -e "  ${YELLOW}<host>${NC}               The hostname to configure (e.g., example.localhost)."
    echo -e "  ${YELLOW}<root>${NC}               The root directory for the host (required for add-host command)."
}

# Parse active hosts from the Caddyfile
list_active_hosts() {
    # Extract hosts that end with .localhost
    grep -E '^[a-zA-Z0-9.-]+\.localhost {' "$CADDYFILE_PATH" | sed -E 's/ \{.*//'
}

# Check and add missing hosts to the hosts file
ensure_hosts_in_hosts_file() {
    local missing_hosts=()
    local hosts=$(list_active_hosts)
    local added_count=0

    # Check if we have root privileges
    if [[ "$EUID" -ne 0 ]]; then
        echo -e "${RED}This operation requires root privileges.${NC}"
        echo -e "Please run one of these commands:"
        echo -e "  ${YELLOW}sudo $0 localhost ensure-hosts${NC}"
        echo -e "  ${YELLOW}sudo $(realpath $0) localhost ensure-hosts${NC}"
        return 1
    fi

    echo -e "\n${YELLOW}Checking hosts in Caddyfile...${NC}"
    
    for host in $hosts; do
        if ! grep -q "$host" "$HOSTS_FILE"; then
            echo -e "${YELLOW}Adding missing host:${NC} $host"
            echo "127.0.0.1 $host" >> "$HOSTS_FILE"
            echo "::1 $host" >> "$HOSTS_FILE"
            missing_hosts+=("$host")
            added_count=$((added_count + 1))
        fi
    done

    if [ $added_count -eq 0 ]; then
        echo -e "${GREEN}All hosts from Caddyfile are already present in the hosts file.${NC}"
    else
        echo -e "\n${GREEN}Successfully added $added_count hosts to $HOSTS_FILE:${NC}"
        printf '%s\n' "${missing_hosts[@]}"
    fi
}

list_hosts() {
    # Print the table header
    printf "\n${THEME}${BOLD}┌────────────────────────────────┬──────────────────────────────────────────┐${NC}\n"
    printf "${THEME}${BOLD}│ %-30s │ %-40s │${NC}\n" "HOST" "ROOT DIRECTORY"
    printf "${THEME}${BOLD}├────────────────────────────────┼──────────────────────────────────────────┤${NC}\n"

    # Initialize variables
    local host=""
    local root=""

    # Read the Caddyfile and extract host and root directory information
    while IFS= read -r line; do
        if [[ $line =~ ^[a-zA-Z0-9.-]+ ]]; then
            host=$(echo $line | tr -d '{')
        elif [[ $line =~ ^[[:space:]]*root[[:space:]]+ ]]; then
            root=$(echo $line | awk '{print $2}')
            printf "${THEME}${BOLD}│${NC} %-30s ${THEME}${BOLD}│${NC} %-40s ${THEME}${BOLD}│${NC}\n" "$host" "$root"
            host=""
            root=""
        fi
    done < "$CADDYFILE_PATH"

    # Print the table footer
    printf "${THEME}${BOLD}└────────────────────────────────┴──────────────────────────────────────────┘${NC}\n"
}

# Main function to handle host modifications
execute_hosts_action() {
    local action=$1
    local host=$2

    case $action in
        add-host)
            add_host_to_hosts_file "$host"
            ;;
        remove-host)
            remove_host_from_hosts_file "$host"
            ;;
        list-hosts)
            list_hosts
            ;;
        ensure-hosts)
            ensure_hosts_in_hosts_file
            ;;
        *)
            echo -e "${RED}Invalid action for hosts file: $action${NC}"
            return 1
            ;;
    esac
}

execute_caddyfile() {
    local action=$1
    local host=$2
    local root=$3

    case $action in
        add-host)
            if [ -z "$root" ]; then
                echo -e "${RED}Root directory is required for adding a host.${NC}"
                show_caddyfile_help
                return 1
            fi

            # Check if the host already exists
            if grep -q "$host" "$CADDYFILE_PATH"; then
                echo -e "${RED}Host $host already exists in the Caddyfile.${NC}"
                return 1
            fi

            # Add new host configuration
            cat <<EOF >> "$CADDYFILE_PATH"

$host {
    import global
    root  $root

    # Add reverse proxy for Vite
    handle /@vite/* {
        reverse_proxy node_devner:5173
    }
}
EOF
            echo -e "${GREEN}Host $host added successfully.${NC}"
            execute_hosts_action add-host "$host"
            make reload -C ${SRCS_DIR}
            return 0
            ;;
        remove-host)
            # Remove the host configuration
            if ! grep -q "$host" "$CADDYFILE_PATH"; then
                echo -e "${RED}Host $host not found in the Caddyfile.${NC}"
                return 1
            fi

            # Create a temporary file
            local temp_file=$(mktemp)
            
            # Use awk to remove the block
            awk -v host="$host" '
                BEGIN { skip = 0; }
                $0 ~ "^" host " {" { skip = 1; next; }
                skip && /^}/ { skip = 0; next; }
                !skip { print; }
            ' "$CADDYFILE_PATH" > "$temp_file"
            
            # Replace the original file with the temporary file
            mv "$temp_file" "$CADDYFILE_PATH"
            
            # Check if the host was removed successfully
            if grep -q "$host" "$CADDYFILE_PATH"; then
                echo -e "${RED}Error: Failed to remove host $host.${NC}\n"
                return 1
            fi

            echo -e "${GREEN}Host $host removed successfully.${NC}\n"
            execute_hosts_action remove-host "$host"
            make reload -C ${SRCS_DIR}

            return 0
            ;;
        list-hosts)
            list_hosts
            ;;
        *)
            echo -e "${RED}Invalid action: $action${NC}"
            show_caddyfile_help
            return 1
            ;;
    esac
}

trust_caddy_root_ca() {
    # Define variables
    CONTAINER_NAME="frankenphp_devner"
    CERT_PATH="./caddy-root.crt"
    CONTAINER_CERT_PATH="/data/caddy/pki/authorities/local/root.crt"

    echo -e "\n${YELLOW}Extracting Root CA from Caddy container...${NC}"
    
    # Copy the Root CA certificate from the container to the host
    docker cp "${CONTAINER_NAME}:${CONTAINER_CERT_PATH}" "${CERT_PATH}" 2>/dev/null
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}Error:${NC} Failed to copy certificate from the container."
        exit 1
    fi
    
    echo -e "${GREEN}Certificate successfully copied to ${CERT_PATH}.${NC}"

    # Detect OS and install certificate
    OS="$(uname -s)"
    case "${OS}" in
        Darwin) # macOS
            echo -e "${YELLOW}Installing certificate on macOS...${NC}"
            sudo security add-trusted-cert -d -r trustRoot -k /Library/Keychains/System.keychain "${CERT_PATH}"
            ;;
        Linux) # Linux
            echo -e "${YELLOW}Installing certificate on Linux...${NC}"
            sudo cp "${CERT_PATH}" /usr/local/share/ca-certificates/caddy-root.crt
            sudo update-ca-certificates
            ;;
        MINGW*|CYGWIN*|MSYS*|Windows_NT) # Windows
            echo -e "${YELLOW}For Windows:${NC} Open 'Manage User Certificates' and import '${CERT_PATH}' into 'Trusted Root Certification Authorities'."
            explorer.exe .
            ;;
        *)
            echo -e "${RED}Unsupported operating system. Please install the certificate manually.${NC}"
            exit 1
            ;;
    esac

    echo -e "${GREEN}Root CA successfully trusted on your system.${NC}"
}

