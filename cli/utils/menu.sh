#!/bin/bash

show_menu() {
    echo -e "\n${BOLD}Select a command to execute:${NC}\n"
    
    # Define commands in categories
    general_commands=(
        "1) up"
        "2) down"
        "3) stop"
        "4) rebuild"
        "5) delete"
        "6) nocache"
        "7) reload"
    )
    
    db_commands=(
        "8) new-mysql"
        "9) remove-mysql"
        "10) new-postgres"
        "11) remove-postgres"
    )
    
    quick_access_commands=(
        "12) postgres"
        "13) mysql"
        "14) node"
        "15) frankenphp"
    )
    
    project_commands=(
        "16) new"
        "17) remove"
    )
    
    other_commands=(
        "18) ps"
        "19) alias"
        "20) code"
        "21) wp"
    )
    
    caddyfile_commands=(
        "22) add-host"
        "23) remove-host"
        "24) list-hosts"
        "25) ensure-hosts"
    )
    
    help_and_quit_commands=(
        "c) credit"
        "h) help"
        "q) quit"
    )
    
    # Function to print commands in two columns
    print_commands() {
        local commands=("$@")
        local num_commands=${#commands[@]}
        local half=$(( (num_commands + 1) / 2 ))

        for ((i = 0; i < half; i++)); do
            cmd1="${commands[$i]}"
            cmd2=""
            if [ $((i + half)) -lt $num_commands ]; then
                cmd2="${commands[$((i + half))]}"
            fi
            printf "${THEME}${BOLD}│${NC} %-25s %-25s ${THEME}${BOLD}│${NC}\n" "$cmd1" "$cmd2"
        done
    }

    # Print the menu
    printf "${THEME}${BOLD}┌─────────────────────────────────────────────────────┐${NC}\n"
    printf "${THEME}${BOLD}│ %-25s %-25s │${NC}\n" "General Commands" ""
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    print_commands "${general_commands[@]}"
    
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    printf "${THEME}${BOLD}│ %-25s %-25s │${NC}\n" "Database Commands" ""
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    print_commands "${db_commands[@]}"
    
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    printf "${THEME}${BOLD}│ %-25s %-25s │${NC}\n" "Quick Access Commands" ""
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    print_commands "${quick_access_commands[@]}"
    
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    printf "${THEME}${BOLD}│ %-25s %-25s │${NC}\n" "Project Commands" ""
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    print_commands "${project_commands[@]}"
    
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    printf "${THEME}${BOLD}│ %-25s %-25s │${NC}\n" "Other Commands" ""
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    print_commands "${other_commands[@]}"
    
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    printf "${THEME}${BOLD}│ %-25s %-25s │${NC}\n" "Caddyfile Commands" ""
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    print_commands "${caddyfile_commands[@]}"
    
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    printf "${THEME}${BOLD}│ %-25s %-25s │${NC}\n" "Help and Quit" ""
    printf "${THEME}${BOLD}├─────────────────────────────────────────────────────┤${NC}\n"
    print_commands "${help_and_quit_commands[@]}"
    printf "${THEME}${BOLD}└─────────────────────────────────────────────────────┘${NC}\n"
}

handle_menu_selection() {
    if [ $STATUS -eq 1 ]; then
        echo -e "${RED}Invalid selection!${NC}"
        STATUS=0
    fi

    read -p "Enter your choice [1-25] or command name (q to quit): " choice

    # Map numbers to command names
    case $choice in
        1 | "up") clear; execute_general_command up ;;
        2 | "down") clear; execute_general_command down ;;
        3 | "stop") clear; execute_general_command stop ;;
        4 | "rebuild") clear; execute_general_command rebuild ;;
        5 | "delete") clear; execute_general_command delete ;;
        6 | "nocache") clear; execute_general_command nocache ;;
        7 | "reload") clear; execute_general_command reload ;;
        8 | "new-mysql")
            clear
            read -p "Enter database name: " db_name
            read -p "Enter database user: " db_user
            read -sp "Enter database password: " db_pass
            echo
            execute_new_db_command new-mysql "$db_name" "$db_user" "$db_pass"
            read -p "Press [Enter] key to continue..."
            ;;
        9 | "remove-mysql")
            clear
            read -p "Enter database name: " db_name
            read -p "Enter database user: " db_user
            echo
            execute_remove_db_command remove-mysql "" "$db_name" "$db_user"
            read -p "Press [Enter] key to continue..."
            ;;
        10 | "new-postgres")
            clear
            read -p "Enter database name: " db_name
            read -p "Enter database user: " db_user
            read -sp "Enter database password: " db_pass
            echo
            execute_new_db_command new-postgres "$db_name" "$db_user" "$db_pass"
            read -p "Press [Enter] key to continue..."
            ;;
        11 | "remove-postgres")
            clear
            read -p "Enter database name: " db_name
            read -p "Enter database user: " db_user
            echo
            execute_remove_db_command remove-postgres "" "$db_name" "$db_user"
            read -p "Press [Enter] key to continue..."
            ;;
        12 | "postgres") clear; execute_db_command postgres ;;
        13 | "mysql") clear; execute_db_command mysql ;;
        14 | "node") clear; execute_quick_access_command node ;;
        15 | "frankenphp") clear; execute_quick_access_command frankenphp ;;
        16 | "new")
            clear
            read -p "Enter project type (laravel/wp): " proj_type
            read -p "Enter project name: " proj_name
            read -p "Enter database type (mysql/postgres): " db_type
            execute_project_command new "$proj_type" "$proj_name" "$db_type"
            read -p "Press [Enter] key to continue..."
            ;;
        17 | "remove")
            clear
            read -p "Enter project name: " proj_name
            read -p "Enter database type (mysql/postgres): " db_type
            execute_project_command remove "" "$proj_name" "$db_type"
            read -p "Press [Enter] key to continue..."
            ;;
        18 | "ps") clear; execute_other_command ps; read -p "Press [Enter] key to continue..." ;;
        19 | "alias") clear; execute_other_command alias; read -p "Press [Enter] key to continue..." ;;
        20 | "code") clear; execute_other_command code ;;
        21 | "wp") 
            clear
            read -p "Enter command <project_name> <args>: " -a args
            execute_wp_cli_menu "${args[@]}"
            read -p "Press [Enter] key to continue..."
            ;;
        22 | "add-host")
            clear
            read -p "Enter host name: " host
            read -p "Enter root directory: " root
            execute_caddyfile add-host "$host" "$root"
            read -p "Press [Enter] key to continue..."
            ;;
        23 | "remove-host")
            clear
            read -p "Enter host name: " host
            execute_caddyfile remove-host "$host"
            read -p "Press [Enter] key to continue..."
            ;;
        24 | "list-hosts") clear; execute_caddyfile list-hosts; read -p "Press [Enter] key to continue..." ;;
        25 | "ensure-hosts")
            clear
            execute_caddyfile ensure-hosts
            read -p "Press [Enter] key to continue..."
            ;;
        c | "credit") clear; handle_credit; echo ; read -p "Press [Enter] key to continue..." ;;
        h | "help") clear; handle_help_selection ;;
        q | "quit") clear; exit 0 ;;
        *) clear; STATUS=1 ;;
    esac

    clear
}
