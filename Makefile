# COLORS
GREEN		= \033[1;32m
RED 		= \033[1;31m
ORANGE		= \033[1;33m
CYAN		= \033[1;36m
RESET		= \033[0m

# FOLDER
SRCS_DIR	= ./
DOCKER_DIR	= ${SRCS_DIR}docker-compose.yml

# VARIABLES
ENV_FILE	= ${SRCS_DIR}.env

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

nocache:
	@echo "${GREEN}Rebuilding containers...${RESET}"
	@${DOCKER} build --no-cache
	@${DOCKER} up -d --remove-orphans

delete:
	@echo "${RED}Deleting containers...${RESET}"
	@${DOCKER} down -v --remove-orphans

node:
	@echo "${GREEN}Entering node container...${RESET}"
	@${DOCKER} exec node bash

frankenphp:
	@echo "${GREEN}Entering frankenphp container...${RESET}"
	@${DOCKER} exec frankenphp bash

mysql8:
	@echo "${GREEN}Entering mysql8 container...${RESET}"
	@${DOCKER} exec mysql_8 bash

mysql5:
	@echo "${GREEN}Entering mysql5 container...${RESET}"
	@${DOCKER} exec mysql5 bash

nuxt:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash

gui-dev:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "yarn dev"

gui-build:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "yarn build"

gui:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "yarn start"

.PHONY: all start up down stop rebuild delete