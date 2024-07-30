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

    # Capture the output of the make command
    local make_output
    make_output=$(make $command -C ${SRCS_DIR} database_name=$dbname database_user=$dbuser database_password=$dbpass 2>&1)

    # Check the output from the make command to determine if the database or user was created
    if grep -q "already exist" <<< "$make_output"; then
        return 1
    elif grep -q "created" <<< "$make_output"; then
        return 0
    else
        return 2
    fi
}

execute_remove_db_command() {
    local command=$1
    local dbname=$2
    local dbuser=$3

    if [ -z "$dbname" ] || [ -z "$dbuser" ]; then
        echo -e "${RED}Missing arguments.${NC} Usage: $0 $command <database_name> <database_user>"
        return 1
    fi

    # Capture the output of the make command
    local make_output
    make_output=$(make $command -C ${SRCS_DIR} database_name=$dbname database_user=$dbuser 2>&1)

    # Check the output from the make command to determine if the database or user was removed
    if grep -q "does not exist" <<< "$make_output"; then
        return 1
    elif grep -q "removed" <<< "$make_output"; then
        return 0
    else
        return 2
    fi
}

execute_db_user_check() {
    local command=$1
    local dbname=$2
    local dbuser=$3

    if [ -z "$dbname" ] || [ -z "$dbuser" ]; then
        echo -e "${RED}Missing arguments.${NC} Usage: $0 $command <database_name> <database_user>"
        return 1
    fi

    # Capture the output of the make command
    local make_output
    make_output=$(make $command -C ${SRCS_DIR} database_name=$dbname database_user=$dbuser 2>&1)

   # Check the output from the make command to determine if the database or user exists
    if grep -q "does not exist" <<< "$make_output"; then
        return 1
    elif grep -q "exists" <<< "$make_output"; then
        return 0
    else
        return 2
    fi
}
