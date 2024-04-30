#!/bin/bash

# Color Codes
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# FOLDER
SRCS_DIR=$(cd $(dirname $0); pwd)
PROJECTS_DIR="${SRCS_DIR}/projects"

show_help() {
    echo -e "Usage: ${YELLOW}$0 <command>${NC}\n"
    echo -e "Run these commands from the root of the cloned repository:\n"
    
    echo -e "${BLUE}General Commands:${NC}"
    echo -e "${YELLOW}- up${NC} - Start the development environment."
    echo -e "${YELLOW}- down${NC} - Stop and remove containers, networks, images, and volumes."
    echo -e "${YELLOW}- stop${NC} - Stop all containers without removing them."
    echo -e "${YELLOW}- rebuild${NC} - Rebuild the Docker containers."
    echo -e "${YELLOW}- delete${NC} - Remove all containers and volumes.\n"
    
    echo -e "${BLUE}Database Access Commands:${NC}"
    echo -e "${YELLOW}- mysql8${NC} - Access the MySQL 8 container."
    echo -e "${YELLOW}- mysql5${NC} - Access the MySQL 5 container.\n"
    
    echo -e "${BLUE}Node and PHP Version-specific Commands:${NC}"
    echo -e "${YELLOW}- node${NC} - Access the Node container (uses NVM for multiple Node.js versions)."
    echo -e "${YELLOW}- php82${NC} - Access the PHP 8.2 container."
    echo -e "${YELLOW}- php81${NC} - Access the PHP 8.1 container."
    echo -e "${YELLOW}- php8${NC} - Access the PHP 8.0 container."
    echo -e "${YELLOW}- php74${NC} - Access the PHP 7.4 container.\n"
    
    echo -e "${BLUE}Additional Commands:${NC}"
    echo -e "${YELLOW}- host${NC} - Add a new host entry to /etc/hosts and create a new virtual host."
    echo -e "${YELLOW}- reload${NC} - Reload Nginx configuration without restarting the container."
    echo -e "${YELLOW}- code${NC} - Open a project in Visual Studio Code by listing available projects."
    echo -e "${YELLOW}- add-alias${NC} - Add an alias to .bashrc and .zshrc for easy script execution."
    echo -e "${YELLOW}- ps${NC} - Check if the devner container is running and start it if not.\n"

    echo -e "${BLUE}New projectCommands:${NC}"
    echo -e "${YELLOW}- new <laravel/wp> <project name>${NC} - Create a new Laravel or WordPress project.\n"
}

# Check if at least one argument is provided
if [ $# -eq 0 ]; then
    show_help
    exit 1
fi

command=$1

case $command in
    up|down|stop|rebuild|delete|nginx|mysql8|mysql5|node|php82|php81|php8|php74|host|reload)
        echo -e "${GREEN}Executing command: ${command}${NC}"
        make ${command} -C ${SRCS_DIR}
        ;;
    ps)
        # Check if the devner container is running
        if docker ps | grep -q "devner"; then
            echo -e "Devner container is ${GREEN}running${NC}."
        else
            echo -e "Devner container is not ${RED}running${NC}."

            read -p "Do you want to start the container? (y/n) " start_container
            if [ "$start_container" = "y" ]; then
                make up -C ${SRCS_DIR}
            fi
        fi
        ;;
    code)
        echo -e "${BLUE}Available Projects:${NC}"
        # List all directories (projects) in the PROJECTS_DIR and add them to an array
        projects=($(ls -d ${PROJECTS_DIR}/*/))
        for i in "${!projects[@]}"; do
            project_name=$(basename "${projects[$i]}")
            echo -e "${YELLOW}$((i+1)))${NC} ${project_name}"
        done

        # Prompt user for a choice
        read -p "#? " project_choice

        # Validate the input is a number and within the range of available projects
        if [[ ! $project_choice =~ ^[0-9]+$ ]] || [ "$project_choice" -lt 1 ] || [ "$project_choice" -gt "${#projects[@]}" ]; then
            echo -e "${RED}Invalid selection. Please enter a number from the list.${NC}"
            exit 1
        fi

        # Adjusting the index to 0-based for array access
        project_index=$((project_choice-1))
        project_path="${projects[$project_index]}"

        # Removing the trailing slash from the project_path
        project_path="${project_path%/}"

        project_name=$(basename "${project_path}")
        echo -e "${GREEN}Opening ${project_name} in VSCode...${NC}"

        # Start docker container if not running
        if ! docker ps | grep -q "devner"; then
            make up -C ${SRCS_DIR}
        fi

        code "${project_path}"
        ;;
    add-alias)
        # Define the alias command
        alias_command="alias devner='bash ${SRCS_DIR}/devner.sh'"

        # Add to .bashrc if it exists
        if [ -f "$HOME/.bashrc" ]; then
            if ! grep -q "alias devner=" "$HOME/.bashrc"; then
                echo "$alias_command" >> "$HOME/.bashrc"
                echo -e "Alias added to ${GREEN}.bashrc${NC}. Restart your terminal or source .bashrc to use it."
            else
                echo -e "Alias ${YELLOW}devner${NC} already exists in .bashrc."
            fi
        fi

        # Add to .zshrc if it exists
        if [ -f "$HOME/.zshrc" ]; then
            if ! grep -q "alias devner=" "$HOME/.zshrc"; then
                echo "$alias_command" >> "$HOME/.zshrc"
                echo -e "Alias added to ${GREEN}.zshrc${NC}. Restart your terminal or source .zshrc to use it."
            else
                echo -e "Alias ${YELLOW}devner${NC} already exists in .zshrc."
            fi
        fi
        ;;
    new)
        if [ $# -lt 3 ]; then
            echo -e "${RED}Please provide the type of project (laravel/wp) and the project name.${NC}"
            exit 1
        fi

        project_type=$2
        project_name=$3

        case $project_type in
            laravel)
                make new-laravel -C ${SRCS_DIR} project_name=${project_name}
                ;;
            wp)
                make new-wp -C ${SRCS_DIR} project_name=${project_name}
                ;;
            *)
                echo -e "${RED}Invalid project type${NC}: ${project_type}"
                echo -e "Please provide either ${YELLOW}laravel${NC} or ${YELLOW}wp${NC} as the project type."
                exit 1
                ;;
        esac
        ;;
    help)
        show_help
        ;;
    *)
        echo -e "\n${RED}Invalid command${NC}: ${command}\n"
        echo -e "Run ${YELLOW}devner help${NC} to see the list of available commands."
        exit 1
        ;;

esac
