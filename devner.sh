#!/bin/bash

# FOLDER
SRCS_DIR=$(cd $(dirname $0); pwd)
PROJECTS_DIR="${SRCS_DIR}/projects"
STATUS=0

# Source style scripts
source "${SRCS_DIR}/cli/utils/style.sh"
source "${SRCS_DIR}/cli/utils/helpers.sh"
source "${SRCS_DIR}/cli/utils/help.sh"
source "${SRCS_DIR}/cli/utils/credit.sh"
source "${SRCS_DIR}/cli/utils/menu.sh"

# Source command scripts
source "${SRCS_DIR}/cli/cmd/general_commands.sh"
source "${SRCS_DIR}/cli/cmd/db_commands.sh"
source "${SRCS_DIR}/cli/cmd/quick_access_commands.sh"
source "${SRCS_DIR}/cli/cmd/other_commands.sh"
source "${SRCS_DIR}/cli/cmd/caddyfile_commands.sh"
source "${SRCS_DIR}/cli/cmd/project_commands.sh"
source "${SRCS_DIR}/cli/cmd/wp_cli_commands.sh"

# Style
THEME=${NC}

# Check if the projects directory exists
if [ ! -d "${PROJECTS_DIR}" ]; then
    mkdir -p "${PROJECTS_DIR}"
fi

# Check if Cadddyfile exists
if [ ! -f "${SRCS_DIR}/services/frankenphp/Caddyfile" ]; then
    cp "${SRCS_DIR}/services/frankenphp/default/Caddyfile.example" "${SRCS_DIR}/services/frankenphp/Caddyfile"
fi

if [ -z "$1" ]; then
    while true; do
        show_menu
        handle_menu_selection
    done
else
    case $1 in
        
        up|down|stop|rebuild|delete|nocache|reload)
            execute_general_command $1
            ;;
        
        remove)
            # Handle: devner remove <project> <db_type>
            if [ -z "$2" ]; then
                echo -e "${RED}Please provide the project name.${NC}"
                exit 1
            fi
            if [ -z "$3" ]; then
                echo -e "${RED}Please provide the database type (mysql/postgres).${NC}"
                exit 1
            fi
            execute_project_command remove unused $2 $3
            ;;
        wp)
            execute_wp_cli $@
            ;;
        node|frankenphp|run)
            execute_quick_access_command $1
            ;;
        ps|alias|code|zed)
            execute_other_command $1 $2
            ;;
        postgres|mysql)
            execute_db_command $1
            ;;
        help|--h)
            show_help
            ;;
        credit|--c)
            handle_credit
            ;;
        new-mysql|new-postgres)
            execute_new_db_command $@
            ;;
        remove-mysql|remove-postgres)
            execute_remove_db_command $@
            ;;
        check-mysql|check-postgres)
            execute_db_user_check $@
            ;;
        new|--n|delete|--d)
            execute_project_command $@
            ;;
        add-host|remove-host|list-hosts)
            execute_caddyfile $@
            ;;
        localhost)
            shift
            case $1 in
                ensure-hosts)
                    execute_hosts_action ensure-hosts
                    ;;
                add-host|remove-host|list-hosts)
                    execute_hosts_action $@
                    ;;
                *)
                    echo -e "${RED}Invalid action for localhost${NC}: $1"
                    exit 1
                    ;;
            esac
            ;;
        trust-ca)
            trust_caddy_root_ca
            ;;
        install)
            execute_other_command install
            ;;
        open)
            make dev -C ${SRCS_DIR}
            ;;
        *)
            echo -e "\n${RED}Invalid command${NC}: $1\n"
            echo -e "Run ${YELLOW}$0 help${NC} to see the list of available commands."
            exit 1
            ;;
    esac
fi