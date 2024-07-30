#!/bin/bash

# Show help for wp-cli usage
show_wp_cli_help() {
    echo -e "Usage: ${BOLD}$0 wp-cli <project_name> <wp_args>\n"
    echo -e "Commands:"
    echo -e "  ${YELLOW}wp-cli <project_name> <wp_args>${NC} Run wp commands"
    echo -e "\nArguments:"
    echo -e "  ${YELLOW}<project_name>${NC}          The project to target"
    echo -e "  ${YELLOW}<wp_args>${NC}               The args of the command."
}

# Function to run wp-cli commands
execute_wp_cli() {
    shift
    local wp_dir=$1
    shift
    local wp_args=$@

    # Show help if missing args or project name
    if [ -z "$wp_dir" ] || [ -z "$wp_args" ]; then
        show_wp_cli_help
        return 1
    fi

    make wp-cli -C ${SRCS_DIR} project_name="$wp_dir" wp_args="$wp_args"
}

# Function to run wp-cli commands
execute_wp_cli_menu() {
    local wp_dir=$1
    shift
    local wp_args=$@

    # Show help if missing args or project name
    if [ -z "$wp_dir" ] || [ -z "$wp_args" ]; then
        show_wp_cli_help
        return 1
    fi

    make wp-cli -C ${SRCS_DIR} project_name="$wp_dir" wp_args="$wp_args"
}