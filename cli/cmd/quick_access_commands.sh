#!/bin/bash

execute_quick_access_command() {
    local command=$1
    echo -e "${GREEN}Accessing container: ${command}${NC}"
    make $command -C ${SRCS_DIR}
}