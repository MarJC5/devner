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

mysql:
	@echo "${GREEN}Entering mysql container...${RESET}"
	@${DOCKER} exec mysql bash

postgres:
	@echo "${GREEN}Entering postgres container...${RESET}"
	@${DOCKER} exec postgres bash

nuxt:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash

gui-clean:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "npx nuxi cleanup"

gui-dev:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "yarn dev"

gui-build:
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "yarn build"

gui: up
	@echo "${GREEN}Entering nuxt container...${RESET}"
	@${DOCKER} exec gui bash -c "yarn start"

new-wp:
	@echo "${GREEN}Creating new wordpress project...${RESET}"
	@$(DOCKER) exec frankenphp bash -c "wp core download --path=${project_name} --locale=fr_FR --allow-root"

new-laravel:
	@echo "${GREEN}Creating new laravel project...${RESET}"
	@$(DOCKER) exec frankenphp bash -c "composer create-project --prefer-dist laravel/laravel ${project_name}"

.PHONY: all start up down stop rebuild delete