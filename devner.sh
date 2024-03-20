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

# Check if at least one argument is provided
if [ $# -eq 0 ]; then
    echo -e "${YELLOW}Usage: $0 <command>${NC}"
    echo -e "Available commands: ${GREEN}up, down, stop, rebuild, delete, nginx, mysql8, mysql5, node, php82, php81, php8, php74, host, reload, code, add-alias${NC}"
    exit 1
fi

command=$1

case $command in
    up|down|stop|rebuild|delete|nginx|mysql8|mysql5|node|php82|php81|php8|php74|host|reload)
        echo -e "${GREEN}Executing command: ${command}${NC}"
        make ${command} -C ${SRCS_DIR}
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
    *)
        echo -e "${RED}Invalid command: ${command}${NC}"
        echo -e "Available commands: ${GREEN}up, down, stop, rebuild, delete, nginx, mysql8, mysql5, node, php82, php81, php8, php74, host, reload, code, add-alias${NC}"
        exit 1
        ;;
esac
