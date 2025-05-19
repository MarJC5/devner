#!/bin/bash

check_ps() {
    if docker ps | grep -q "devner"; then
        echo -e "Devner container is ${GREEN}running${NC}."

        # Capture the formatted docker ps output
        docker_output=$(docker ps --filter "name=devner" --format "{{.Names}}\t{{.Status}}\t{{.RunningFor}}\t{{.Ports}}")

        # Print the table header
        printf "\n${THEME}${BOLD}┌────────────────────┬──────────────────────────────┬───────────────────┬─────────────────────────┐${NC}\n"
        printf "${THEME}${BOLD}│ %-18s │ %-28s │ %-17s │ %-23s │${NC}\n" "NAMES" "STATUS" "RUNNING FOR" "PORTS"
        printf "${THEME}${BOLD}├────────────────────┼──────────────────────────────┼───────────────────┼─────────────────────────┤${NC}\n"

        # Print each line of the docker output in a formatted table row
        echo "$docker_output" | while IFS=$'\t' read -r name status runningfor ports; do
            wrapped_ports=$(wrap_text "$ports" 23) # Wrap ports to fit within 23 characters
            printf "${THEME}${BOLD}│${NC} %-18s ${THEME}${BOLD}│${NC} %-28s ${THEME}${BOLD}│${NC} %-17s ${THEME}${BOLD}│${NC} %-23s ${THEME}${BOLD}│${NC}\n" "$name" "$status" "$runningfor" "$(echo "$wrapped_ports" | head -n1)"
            echo "$wrapped_ports" | tail -n +2 | while read -r line; do
                printf "${THEME}${BOLD}│${NC} %-18s ${THEME}${BOLD}│${NC} %-28s ${THEME}${BOLD}│${NC} %-17s ${THEME}${BOLD}│${NC} %-23s ${THEME}${BOLD}│${NC}\n" "" "" "" "$line"
            done
        done

        # Print the table footer
        printf "${THEME}${BOLD}└────────────────────┴──────────────────────────────┴───────────────────┴─────────────────────────┘${NC}\n"
    else
        echo -e "Devner container is not ${RED}running${NC}."

        read -p "Do you want to start the container? (y/n) " start_container
        if [ "$start_container" = "y" ]; then
            make up -C ${SRCS_DIR}
        fi
    fi
}

open_code() {
    # Check if $1 is not empty to open a specific project or show a list of projects
    if [ -n "$1" ]; then
        local project_name=$1
        local project_path="${PROJECTS_DIR}/${project_name}"
        if [ -d "$project_path" ]; then
            echo -e "${GREEN}Opening ${project_name} in VSCode...${NC}"
            code "${project_path}"
        else
            echo -e "${RED}Project ${project_name} not found.${NC}"
        fi
        return
    fi

    echo -e "\n${UNDERLINE}${BOLD}Available Projects:${NC}\n"
    projects=($(ls -d ${PROJECTS_DIR}/*/))
    local half=$(((${#projects[@]} + 1) / 2))

    # Print header
    printf "${THEME}${BOLD}┌─────┬─────────────────────────┬─────┬─────────────────────────┐${NC}\n"
    printf "${THEME}${BOLD}│ %-3s │ %-23s │ %-3s │ %-23s │${NC}\n" "No." "Project" "No." "Project"
    printf "${THEME}${BOLD}├─────┼─────────────────────────┼─────┼─────────────────────────┤${NC}\n"

    # Print projects in two columns
    for ((i=0; i<half; i++)); do
        local index1=$((i+1))
        local index2=$((i+1+half))
        local project_name1=$(basename "${projects[$i]}")
        local project_name2=""

        if [ $index2 -le ${#projects[@]} ]; then
            project_name2=$(basename "${projects[$((i+half))]}")
            printf "${THEME}${BOLD}│${NC} %-3s ${THEME}${BOLD}│${NC} %-23s ${THEME}${BOLD}│${NC} %-3s ${THEME}${BOLD}│${NC} %-23s ${THEME}${BOLD}│${NC}\n" $index1 "$(wrap_text "$project_name1" 23)" $index2 "$(wrap_text "$project_name2" 23)"
        else
            printf "${THEME}${BOLD}│${NC} %-3s ${THEME}${BOLD}│${NC} %-23s ${THEME}${BOLD}│${NC} %-3s ${THEME}${BOLD}│${NC} %-23s ${THEME}${BOLD}│${NC}\n" $index1 "$(wrap_text "$project_name1" 23)" "" ""
        fi
    done

    # Print footer
    printf "${THEME}${BOLD}└─────┴─────────────────────────┴─────┴─────────────────────────┘${NC}\n"

    while true; do
        echo -e "\nEnter the number of the project you want to open (q to quit):"
        read -p "#? " project_choice

        # q to quit
        if [ "$project_choice" = "q" ]; then
            return 1
        fi

        if ! [[ "$project_choice" =~ ^[0-9]+$ ]]; then
            echo -e "\n${RED}Invalid selection. Please enter a number from the list or 'q' to quit.${NC}"
        else
            local project_index=$((project_choice-1))
            if [ "$project_index" -ge 0 ] && [ "$project_index" -lt "${#projects[@]}" ]; then
                local project_path="${projects[$project_index]}"
                project_path="${project_path%/}"
                local project_name=$(basename "${project_path}")
                echo -e "${GREEN}Opening ${project_name} in VSCode...${NC}"

                if ! docker ps | grep -q "devner"; then
                    make up -C ${SRCS_DIR}
                fi

                code "${project_path}"
                break
            else
                echo -e "${RED}Invalid selection. Please enter a number from the list or 'q' to quit.${NC}"
            fi
        fi
    done
}

add_alias() {
    local alias_command="alias devner='bash ${SRCS_DIR}/devner.sh'"

    if [ -f "$HOME/.bashrc" ]; then
        if ! grep -q "alias devner=" "$HOME/.bashrc"; then
            echo "$alias_command" >> "$HOME/.bashrc"
            echo -e "Alias added to ${GREEN}.bashrc${NC}. Restart your terminal or source .bashrc to use it."
        else
            sed -i '' "s|alias devner=.*|$alias_command|" "$HOME/.bashrc"
            echo -e "Alias updated in ${GREEN}.bashrc${NC}. Restart your terminal or source ~/.bashrc to use it."
        fi
    fi

    if [ -f "$HOME/.zshrc" ]; then
        if ! grep -q "alias devner=" "$HOME/.zshrc"; then
            echo "$alias_command" >> "$HOME/.zshrc"
            echo -e "Alias added to ${GREEN}.zshrc${NC}. Restart your terminal or source .zshrc to use it."
        else
            sed -i '' "s|alias devner=.*|$alias_command|" "$HOME/.zshrc"
            echo -e "Alias updated in ${GREEN}.zshrc${NC}. Restart your terminal or source ~/.zshrc to use it."
        fi
    fi
}

remove_alias() {
    local alias_pattern="alias devner="

    # Function to remove alias from a specific file
    remove_alias_from_file() {
        local file=$1
        if [ -f "$file" ]; then
            if grep -q "$alias_pattern" "$file"; then
                if [[ "$OSTYPE" == "darwin"* ]]; then
                    # macOS requires a backup suffix for -i or an empty string
                    sed -i '' "/$alias_pattern/d" "$file"
                else
                    # Linux and other systems
                    sed -i "/$alias_pattern/d" "$file"
                fi
                echo -e "Alias removed from ${GREEN}${file}${NC}."
            else
                echo -e "No alias found in ${GREEN}${file}${NC}."
            fi
        fi
    }

    # Remove alias from bashrc and zshrc if they exist
    remove_alias_from_file "$HOME/.bashrc"
    remove_alias_from_file "$HOME/.zshrc"
}

install_to_system() {
    local target_dir="/usr/local/bin"
    local script_path="$(realpath "$0")"
    local script_name="devner"

    # Check if we have root privileges
    if [[ "$EUID" -ne 0 ]]; then
        echo -e "${RED}This operation requires root privileges.${NC}"
        echo -e "Please run: ${YELLOW}sudo ${script_path} install${NC}"
        return 1
    fi

    # Check if the script is already installed
    if [ -f "${target_dir}/${script_name}" ]; then
        echo -e "${YELLOW}Devner is already installed at ${target_dir}/${script_name}${NC}"
        read -p "Do you want to reinstall it? (y/n) " reinstall
        if [ "$reinstall" != "y" ]; then
            return 0
        fi
    fi

    # Create a wrapper script
    cat > "${target_dir}/${script_name}" << EOF
#!/bin/bash
"${script_path}" "\$@"
EOF

    # Make it executable
    chmod +x "${target_dir}/${script_name}"

    echo -e "${GREEN}Devner has been installed to ${target_dir}/${script_name}${NC}"
    echo -e "You can now run it from anywhere using: ${YELLOW}devner${NC}"
    echo -e "Example commands:"
    echo -e "  ${YELLOW}devner up${NC} - Start the development environment"
    echo -e "  ${YELLOW}sudo devner localhost ensure-hosts${NC} - Ensure all hosts are in your hosts file"
}

execute_other_command() {
    local command=$1
    local argument=$2
    case $command in
        ps)
            check_ps
            ;;
        code)
            open_code $argument
            ;;
        alias)
            case $argument in
                add)
                    add_alias
                    ;;
                remove)
                    remove_alias
                    ;;
                *)
                    # Check missing or invalid argument
                    if [ -z "$argument" ]; then
                        echo -e "${RED}Missing argument for alias command${NC}. Please specify 'add' or 'remove'."
                        exit 1
                    fi

                    echo -e "${RED}Invalid argument for alias command${NC}: $argument"
                    exit 1
                    ;;
            esac
            ;;
        path)
            add_to_path
            ;;
        install)
            install_to_system
            ;;
    esac
}
