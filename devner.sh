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

# Style
THEME=${NC}

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
        postgres|mysql)
            execute_db_command $1
            ;;
        new-mysql|new-postgres)
            shift
            execute_new_db_command $1 $@
            ;;
        remove-mysql|remove-postgres)
            shift
            execute_remove_db_command $1 $@
            ;;
        node|frankenphp)
            execute_quick_access_command $1
            ;;
        new|delete)
            shift
            execute_project_command $1 $@
            ;;
        ps|alias|code)
            execute_other_command $1
            ;;
        add-host|remove-host|list-hosts)
            shift
            execute_caddyfile $1 $@
            ;;
        help|--h)
            show_help
            ;;
        credit|--c)
            handle_credit
            ;;
        *)
            echo -e "\n${RED}Invalid command${NC}: $1\n"
            echo -e "Run ${YELLOW}$0 help${NC} to see the list of available commands."
            exit 1
            ;;
    esac
fi