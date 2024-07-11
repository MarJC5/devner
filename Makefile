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
DOCKER		=  docker compose -f ${DOCKER_DIR} --env-file ${ENV_FILE} -p devner

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

nuxt:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash

nuxt-start:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "yarn start"

nuxt-prod:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "yarn build"

.PHONY: all start up down stop rebuild delete