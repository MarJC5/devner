#!/bin/bash

show_help() {
    echo -e "Usage: ${BOLD}$0 <command>${NC}\n"
    echo -e "Run these commands from the root of the cloned repository:\n"
    
    echo -e "${BOLD}General Commands:${NC}"
    echo -e "${YELLOW}- up${NC} - Start the development environment."
    echo -e "${YELLOW}- down${NC} - Stop and remove containers, networks, images, and volumes."
    echo -e "${YELLOW}- stop${NC} - Stop all containers without removing them."
    echo -e "${YELLOW}- rebuild${NC} - Rebuild the Docker containers."
    echo -e "${YELLOW}- delete${NC} - Remove all containers and volumes."
    echo -e "${YELLOW}- nocache${NC} - Remove all containers and volumes and remove cache."
    echo -e "${YELLOW}- reload${NC} - Reload the frankenphp container to update Caddyfile.\n"
    
    echo -e "${BOLD}Database Access Commands:${NC}"
    echo -e "${YELLOW}- new-mysql <database_name> <database_user> <database_password>${NC} - Create a new MySQL database and user."
    echo -e "${YELLOW}- new-postgres <database_name> <database_user> <database_password>${NC} - Create a new PostgreSQL database and user."
    echo -e "${YELLOW}- remove-mysql <database_name> <database_user>${NC} - Remove a MySQL database and user."
    echo -e "${YELLOW}- remove-postgres <database_name> <database_user>${NC} - Remove a PostgreSQL database and user.\n"
    
    echo -e "${BOLD}Quick Access Commands:${NC}"
    echo -e "${YELLOW}- postgres${NC} - Access the PostgreSQL container."
    echo -e "${YELLOW}- mysql${NC} - Access the MySQL 8 container."
    echo -e "${YELLOW}- node${NC} - Access the Node container (uses NVM for multiple Node.js versions)."
    echo -e "${YELLOW}- frankenphp${NC} - Access the Frankenphp container.\n"
    
    echo -e "${BOLD}Project Commands:${NC}"
    echo -e "${YELLOW}- new <laravel/wp> <project name> <mysql/postgres>${NC} - Create a new Laravel or WordPress project."
    echo -e "${YELLOW}- delete <project name> <mysql/postgres>${NC} - Delete a project and its database.\n"
    
    echo -e "${BOLD}Other Commands:${NC}"
    echo -e "${YELLOW}- ps${NC} - Check if the devner container is running."
    echo -e "${YELLOW}- alias${NC} - Add the devner alias to .bashrc or .zshrc."
    echo -e "${YELLOW}- code${NC} - Open a project in VSCode.\n"

    echo -e "${BOLD}Caddyfile Management Commands:${NC}"
    echo -e "${YELLOW}- add-host <host> <root>${NC} - Add a new host to the Caddyfile."
    echo -e "${YELLOW}- remove-host <host>${NC} - Remove an existing host from the Caddyfile."
    echo -e "${YELLOW}- list-hosts${NC} - List all hosts in the Caddyfile.\n"
}



handle_help_selection() {
    clear
    show_help
    while true; do
        read -p "Enter command name (or 'q' to quit help): " command_name

        # if enter or space is pressed without a command name or 'q' is entered, break the loop
        if [ -z "$command_name" ]; then
            break
        fi

        if [ "$command_name" == "q" ]; then
            break
        else
            $0 $command_name
        fi
    done
    clear
}