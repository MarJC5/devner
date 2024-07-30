#!/bin/bash

execute_general_command() {
    local command=$1
    echo -e "${GREEN}Executing command: ${command}${NC}"

    local current_dir=""
    local current_pwd=$(pwd)
    local current_basename=$(basename $current_pwd)

    for project in ${PROJECTS_DIR}/*/; do
        local project_basename=$(basename "$project")
        if [ "$current_basename" = "$project_basename" ]; then
            if [ -f "$project/wp-config.php" ]; then
                current_dir="$(basename "$project")/wp-content/themes"
            else 
                current_dir="$(basename "$project")"
            fi
            break
        fi
    done

    if [ -n "$current_dir" ]; then
        make $command -C ${SRCS_DIR} project_dir=${current_dir}
    else
        make $command -C ${SRCS_DIR}
    fi
}
