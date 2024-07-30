#!/bin/bash

execute_project_command() {
    local command=$1
    local project_type=$2
    local project_name=$3
    local project_db_type=$4
    local host_name="${project_name}"
    local root_dir="/var/www/html"

    # Check args
    if [ -z "$command" ]; then
        echo -e "${RED}Please provide a command.${NC}"
        return 1
    fi

    # Check if the command is valid
    if [ "$command" != "new" ] && [ "$command" != "delete" ]; then
        echo -e "${RED}Invalid command${NC}: ${command}"
        echo -e "Please provide either ${YELLOW}new${NC} or ${YELLOW}delete${NC} as the command."
        return 1
    fi

    # Check if the project type is provided
    if [ "$command" = "new" ] && [ -z "$project_type" ]; then
        echo -e "${RED}Please provide the type of project (laravel/wp).${NC}"
        return 1
    fi

    # Check if the project name is provided
    if [ "$command" = "new" ] && [ -z "$project_name" ]; then
        echo -e "${RED}Please provide the project name.${NC}"
        return 1
    fi

    # Check if the project db type is provided
    if [ "$command" = "new" ] && [ -z "$project_db_type" ]; then
        echo -e "${RED}Please provide the database type (mysql/postgres).${NC}"
        return 1
    fi

    # Check if correct db type is provided
    if [ "$command" = "new" ] && [ "$project_db_type" != "mysql" ] && [ "$project_db_type" != "postgres" ]; then
        echo -e "${RED}Invalid database type${NC}: ${project_db_type}"
        echo -e "Please provide either ${YELLOW}mysql${NC} or ${YELLOW}postgres${NC} as the database type."
        return 1
    fi

     # if not ends with .localhost add it
    if [[ ! $host_name =~ \.localhost$ ]]; then
        # append .localhost to the host name
        host_name="${host_name}.localhost"
    fi

    # if laravel appends /public to the root dir
    if [ "$project_type" = "laravel" ]; then
        root_dir="${root_dir}/${project_name}/public"
    else
        root_dir="${root_dir}/${project_name}"
    fi

    if [ "$command" = "new" ]; then
        if [ $# -lt 3 ]; then
            echo -e "${RED}Please provide the type of project (laravel/wp) and the project name.${NC}"
            return 1
        fi

        # Check if the project already exists
        if [ -d "${SRCS_DIR}/${project_name}" ]; then
            echo -e "${RED}Project already exists.${NC}"
            return 1
        fi

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
                return 1
                ;;
        esac

        if [ "$project_db_type" = "mysql" ]; then
            execute_new_db_command new-mysql ${project_name} ${project_name} ${project_name}
        elif [ "$project_db_type" = "postgres" ]; then
            execute_new_db_command new-postgres ${project_name} ${project_name} ${project_name}
        else
            echo -e "${RED}Invalid database type${NC}: ${project_db_type}"
            echo -e "Please provide either ${YELLOW}mysql${NC} or ${YELLOW}postgres${NC} as the database type."
            return 1
        fi

        execute_caddyfile add-host ${host_name} "/var/www/html/${project_name}"

        echo -e "${GREEN}Project created successfully.${NC}"
    fi

    if [ "$command" = "delete" ]; then
        if [ $# -lt 2 ]; then
            echo -e "${RED}Please provide the project name.${NC}"
            return 1
        fi

        # Check if the project exists
        if [ ! -d "${SRCS_DIR}/${project_name}" ]; then
            echo -e "${RED}Project does not exist.${NC}"
            return 1
        fi

        make remove -C ${SRCS_DIR} project_name=${project_name}

        if [ "$project_db_type" = "mysql" ]; then
            execute_remove_db_command remove-mysql ${project_name} ${project_name}
        elif [ "$project_db_type" = "postgres" ]; then
            execute_remove_db_command remove-postgres ${project_name} ${project_name}
        else
            echo -e "${RED}Invalid database type${NC}: ${project_db_type}"
            echo -e "Please provide either ${YELLOW}mysql${NC} or ${YELLOW}postgres${NC} as the database type."
            return 1
        fi

        execute_caddyfile remove-host ${host_name}

        echo -e "${GREEN}Project deleted successfully.${NC}"
    fi
}
