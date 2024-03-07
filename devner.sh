#!/bin/bash

# FOLDER
SRCS_DIR=$(cd $(dirname $0); pwd)

# Check if at least one argument is provided
if [ $# -eq 0 ]; then
    echo "Usage: $0 <command>"
    echo "Available commands: up, down, stop, rebuild, delete, nginx, mysql8, mysql5, node, php82, php81, php8, php74, host, reload"
    exit 1
fi

command=$1

case $command in
    up)
        make up -C ${SRCS_DIR}
        ;;
    down)
        make down -C ${SRCS_DIR}
        ;;
    stop)
        make stop -C ${SRCS_DIR}
        ;;
    rebuild)
        make rebuild -C ${SRCS_DIR}
        ;;
    delete)
        make delete -C ${SRCS_DIR}
        ;;
    nginx)
        make nginx -C ${SRCS_DIR}
        ;;
    mysql8)
        make mysql8 -C ${SRCS_DIR}
        ;;
    mysql5)
        make mysql5 -C ${SRCS_DIR}
        ;;
    node)
        make node -C ${SRCS_DIR}
        ;;
    php82)
        make php82 -C ${SRCS_DIR}
        ;;
    php81)
        make php81 -C ${SRCS_DIR}
        ;;
    php8)
        make php8 -C ${SRCS_DIR}
        ;;
    php74)
        make php74 -C ${SRCS_DIR}
        ;;
    host)
        make host -C ${SRCS_DIR}
        ;;
    reload)
        make reload -C ${SRCS_DIR}
        ;;
    *)
        echo "Invalid command: $command"
        echo "Available commands: up, down, stop, rebuild, delete, nginx, mysql8, mysql5, node, php82, php81, php8, php74, host, reload"
        exit 1
        ;;
esac
