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
    local alias_command="alias devner='bash \${SRCS_DIR}/devner.sh'"

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

execute_other_command() {
    local command=$1
    case $command in
        ps)
            check_ps
            ;;
        code)
            open_code
            ;;
        alias)
            add_alias
            ;;
    esac
}
