#!/bin/bash

# Path to the Caddyfile
CADDYFILE_PATH="${SRCS_DIR}/services/frankenphp/Caddyfile"

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


execute_caddyfile() {
    local action=$1
    local host=$2
    local root=$3

    case $action in
        add-host)
            if [ -z "$root" ]; then
                echo -e "${RED}Error: Root directory is required for adding a host.${NC}"
                show_caddyfile_help
                return 1
            fi

            # Check if the host already exists
            if grep -q "$host" "$CADDYFILE_PATH"; then
                echo -e "${RED}Error: Host $host already exists in the Caddyfile.${NC}"
                return 1
            fi

            # Add new host configuration
            cat <<EOF >> "$CADDYFILE_PATH"

$host {
    root  $root
    file_server
    php_server
}
EOF
            echo -e "${GREEN}Host $host added successfully.${NC}"
            make reload -C ${SRCS_DIR}
            ;;
        remove-host)
            # Remove the host configuration
            if ! grep -q "$host" "$CADDYFILE_PATH"; then
                echo -e "${RED}Error: Host $host not found in the Caddyfile.${NC}"
                return 1
            fi

            # Use sed to delete the lines between the host declaration and the next blank line
            # Uncoment the following line to generate a backup file
            # sed -i.bak "/$host {/,/^$/d" "$CADDYFILE_PATH"
            echo -e "${GREEN}Host $host removed successfully.${NC}"
            make reload -C ${SRCS_DIR}
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
