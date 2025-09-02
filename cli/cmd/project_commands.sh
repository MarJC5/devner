#!/bin/bash

execute_project_command() {
    local command=$1
    local project_type=$2
    local project_name=$3
    local project_db_type=$4
    local host_name="${project_name}"
    local root_dir="/var/www/html"
    local steps=0
    local total_steps_new=3
    local total_steps_remove=3
    local make_output
    local step_status=()
    local step_messages=()

    # Check args
    if [ -z "$command" ]; then
        echo -e "${RED}Please provide a command.${NC}"
        return 1
    fi

    # Check if the command is valid
    if [ "$command" != "new" ] && [ "$command" != "remove" ]; then
        echo -e "${RED}Invalid command${NC}: ${command}"
        echo -e "Please provide either ${YELLOW}new${NC} or ${YELLOW}remove${NC} as the command."
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
                make_output=$(make new-laravel -C ${SRCS_DIR} project_name=${project_name} 2>&1)
                ;;
            wp)
                make_output=$(make new-wp -C ${SRCS_DIR} project_name=${project_name} 2>&1)
                ;;
            *)
                echo -e "${RED}Invalid project type${NC}: ${project_type}"
                echo -e "Please provide either ${YELLOW}laravel${NC} or ${YELLOW}wp${NC} as the project type."
                return 1
                ;;
        esac

        # Check the output from the make command to determine if the project was created and in
        if grep -q "Error: WordPress files seem to already be present here." <<< "$make_output"; then
            step_status+=("Skipped")
            step_messages+=("WordPress files seem to already be present here.")
        elif grep -q "is not empty" <<< "$make_output"; then
            step_status+=("Skipped")
            step_messages+=("Laravel files seem to already be present here.")
        else
            steps=$((steps + 1))
            step_status+=("Success")
            step_messages+=("Project ${project_type} created.")
        fi

        if [ "$project_db_type" = "mysql" ]; then
            execute_db_user_check check-mysql ${project_name} ${project_name}

            case $? in
                1)
                    execute_new_db_command new-mysql ${project_name} ${project_name} ${project_name}
                    steps=$((steps + 1))
                    step_status+=("Success")
                    step_messages+=("MySQL database and user created.")
                    ;;
                0)
                    step_status+=("Skipped")
                    step_messages+=("MySQL database and user already exist.")
                    ;;
                *)
                    step_status+=("Failed")
                    step_messages+=("Failed to check or create MySQL database and user.")
                    ;;
            esac
        elif [ "$project_db_type" = "postgres" ]; then
            execute_db_user_check check-postgres ${project_name} ${project_name}

            case $? in
                1)
                    execute_new_db_command new-postgres ${project_name} ${project_name} ${project_name}
                    steps=$((steps + 1))
                    step_status+=("Success")
                    step_messages+=("PostgreSQL database and user created.")
                    ;;
                0)
                    step_status+=("Skipped")
                    step_messages+=("PostgreSQL database and user already exist.")
                    ;;
                *)
                    step_status+=("Failed")
                    step_messages+=("Failed to check or create PostgreSQL database and user.")
                    ;;
            esac
        else
            echo -e "${RED}Invalid database type${NC}: ${project_db_type}"
            echo -e "Please provide either ${YELLOW}mysql${NC} or ${YELLOW}postgres${NC} as the database type."
            return 1
        fi

        execute_caddyfile add-host ${host_name} "${root_dir}" > /dev/null

        # check return value of add command
        if [ $? -eq 0 ]; then
            steps=$((steps + 1))
            step_status+=("Success")
            step_messages+=("Host ${host_name} added to Caddyfile.")
        elif [ $? -eq 1 ]; then
            step_status+=("Skipped")
            step_messages+=("Host ${host_name} already exists in Caddyfile.")
        else
            step_status+=("Failed")
            step_messages+=("Failed to add host ${host_name} to Caddyfile.")
        fi

        if [ $steps -eq $total_steps_new ]; then
            echo -e "${GREEN}Project created successfully.${NC}"
        else
            echo -e "${RED}Error creating project, see the summary below.${NC}"
        fi
    fi

    if [ "$command" = "remove" ]; then
        # Since we're calling with: execute_project_command remove unused <project_name> <db_type>
        # The project_name is actually in $3 and db_type in $4 for remove command
        if [ "$2" = "unused" ]; then
            project_name=$3
            project_db_type=$4
        fi
        
        if [ -z "$project_name" ]; then
            echo -e "${RED}Please provide the project name.${NC}"
            return 1
        fi
        
        if [ -z "$project_db_type" ]; then
            echo -e "${RED}Please provide the database type (mysql/postgres).${NC}"
            return 1
        fi
        
        # Regenerate host_name for remove command
        host_name="${project_name}"
        if [[ ! $host_name =~ \.localhost$ ]]; then
            host_name="${host_name}.localhost"
        fi

        # Check if the project exists
        if [ ! -d "${SRCS_DIR}/projects/${project_name}" ]; then
            step_status+=("Failed")
            step_messages+=("Project does not exist.")
        else
            # Remove the project directory
            make remove -C ${SRCS_DIR} project_name=${project_name}
            steps=$((steps + 1))
            step_status+=("Success")
            step_messages+=("Project directory removed.")
        fi

        if [ "$project_db_type" = "mysql" ]; then
            execute_db_user_check check-mysql ${project_name} ${project_name}

            case $? in
                0)
                    execute_remove_db_command remove-mysql ${project_name} ${project_name}
                    steps=$((steps + 1))
                    step_status+=("Success")
                    step_messages+=("MySQL database and user removed.")
                    ;;
                1)
                    step_status+=("Skipped")
                    step_messages+=("MySQL database and user does not exist.")
                    ;;
                *)
                    step_status+=("Failed")
                    step_messages+=("Failed to check or remove MySQL database and user.")
                    ;;
            esac
        elif [ "$project_db_type" = "postgres" ]; then
            execute_db_user_check check-postgres ${project_name} ${project_name}

            case $? in
                0)
                    execute_remove_db_command remove-postgres ${project_name} ${project_name}
                    steps=$((steps + 1))
                    step_status+=("Success")
                    step_messages+=("PostgreSQL database and user removed.")
                    ;;
                1)
                    step_status+=("Skipped")
                    step_messages+=("PostgreSQL database and user does not exist.")
                    ;;
                *)
                    step_status+=("Failed")
                    step_messages+=("Failed to check or remove PostgreSQL database and user.")
                    ;;
            esac
        else
            echo -e "${RED}Invalid database type${NC}: ${project_db_type}"
            echo -e "Please provide either ${YELLOW}mysql${NC} or ${YELLOW}postgres${NC} as the database type."
            return 1
        fi

        execute_caddyfile remove-host ${host_name} > /dev/null

        # check return value of remove command
        if [ $? -eq 0 ]; then
            steps=$((steps + 1))
            step_status+=("Success")
            step_messages+=("Host ${host_name} removed from Caddyfile.")
        elif [ $? -eq 1 ]; then
            step_status+=("Skipped")
            step_messages+=("Host ${host_name} already removed from Caddyfile.")
        else
            step_status+=("Failed")
            step_messages+=("Failed to remove host ${host_name} from Caddyfile.")
        fi

        if [ $steps -eq $total_steps_remove ]; then
            echo -e "${GREEN}Project deleted successfully.${NC}"
        else
            echo -e "${RED}Error deleting project, see the summary below.${NC}"
        fi
    fi

    # Summary of steps
    echo -e "\n${BOLD}${UNDERLINE}Summary:${NC}\n"
    for i in "${!step_status[@]}"; do
        if [ "${step_status[$i]}" = "Success" ]; then
            step_status[$i]="${GREEN}${step_status[$i]}${NC}"
        elif [ "${step_status[$i]}" = "Failed" ]; then
            step_status[$i]="${RED}${step_status[$i]}${NC}"
        elif [ "${step_status[$i]}" = "Skipped" ]; then
            step_status[$i]="${YELLOW}${step_status[$i]}${NC}"
        fi 
        echo -e "[${step_status[$i]}]:${NC} ${step_messages[$i]}"
    done
    echo
}
