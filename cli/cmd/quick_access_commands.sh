#!/bin/bash

execute_quick_access_command() {
    local command=$1
    local current_dir=$(pwd)  # Get the current directory dynamically

    # Extract the last directory name and prefix it with a '/'
    current_dir="/$(basename "$current_dir")"

    # Look into the projects directory if project is wordpress or any other
    # case wordpress => must have wp-config.php and go to themes directory
    # case other => go to root directory

    if [ -f "${PROJECTS_DIR}${current_dir}/wp-config.php" ]; then
        current_dir="${current_dir}/wp-content/themes"
    fi

    # Is project directory exists
    if [ ! -d "${PROJECTS_DIR}${current_dir}" ]; then
        current_dir="/"
    fi

    echo -e "${GREEN}Accessing container: ${command}${NC}"
    make $command -C "${SRCS_DIR}" CURRENT_DIR="${current_dir}"
}