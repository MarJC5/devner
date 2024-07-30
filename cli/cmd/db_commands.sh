#!/bin/bash

execute_db_command() {
    local command=$1
    make $command -C ${SRCS_DIR}
}

execute_new_db_command() {
    local command=$1
    local dbname=$2
    local dbuser=$3
    local dbpass=$4

    if [ -z "$dbname" ] || [ -z "$dbuser" ] || [ -z "$dbpass" ]; then
        echo -e "${RED}Missing arguments.${NC} Usage: $0 $command <database_name> <database_user> <database_password>"
        return 1
    fi

    make $command -C ${SRCS_DIR} database_name=$dbname database_user=$dbuser database_password=$dbpass
}

execute_remove_db_command() {
    local command=$1
    local dbname=$2
    local dbuser=$3

    if [ -z "$dbname" ] || [ -z "$dbuser" ]; then
        echo -e "${RED}Missing arguments.${NC} Usage: $0 $command <database_name> <database_user>"
        return 1
    fi

    make $command -C ${SRCS_DIR} database_name=$dbname database_user=$dbuser
}