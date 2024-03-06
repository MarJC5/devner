# COLORS
GREEN		= \033[1;32m
RED 		= \033[1;31m
ORANGE		= \033[1;33m
CYAN		= \033[1;36m
RESET		= \033[0m

# FOLDER
SRCS_DIR	= ./
DOCKER_DIR	= ${SRCS_DIR}docker-compose.yml


# COMMANDS
DOCKER		=  docker compose -f ${DOCKER_DIR} -p devner

%:
	@:

all: up

start: up

up:
	@echo "${GREEN}Starting containers...${RESET}"
	@${DOCKER} up -d --remove-orphans

down:
	@echo "${RED}Stopping containers...${RESET}"
	@${DOCKER} down

stop:
	@echo "${RED}Stopping containers...${RESET}"
	@${DOCKER} stop

rebuild:
	@echo "${GREEN}Rebuilding containers...${RESET}"
	@${DOCKER} up -d --remove-orphans --build

delete:
	@echo "${RED}Deleting containers...${RESET}"
	@${DOCKER} down -v --remove-orphans

nginx:
	@echo "${GREEN}Running nginx ...${RESET}"
	@${DOCKER} exec nginx sh

mysql8:
	@echo "${GREEN}Running mysql 8 ...${RESET}"
	@${DOCKER} exec mysql_8 bash

mysql5:
	@echo "${GREEN}Running mysql 5.7 ...${RESET}"
	@${DOCKER} exec mysql_5 bash

node:
	@echo "${GREEN}Running nodejs ...${RESET}"
	@${DOCKER} exec node bash

php82:
	@echo "${GREEN}Running php 8.2 ...${RESET}"
	@${DOCKER} exec php_8_2 bash

php81:
	@echo "${GREEN}Running php 8.1 ...${RESET}"
	@${DOCKER} exec php_8_1 bash

php8:
	@echo "${GREEN}Running php 8.0 ...${RESET}"
	@${DOCKER} exec php_8 bash

php74:
	@echo "${GREEN}Running php 7.4 ...${RESET}"
	@${DOCKER} exec php_7_4 bash

host:
	@echo "${GREEN}Adding host...${RESET}"
	@./services/nginx/scripts/host.sh

reload:
	@echo "${GREEN}Reloading nginx...${RESET}"
	@${DOCKER} exec nginx nginx -s reload

.PHONY: all up down stop rebuild delete mysql8 mysql5 node php82 php81 php8 php74