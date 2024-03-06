# Devner

## Description

Boilerplate for working in multiple PH, Node versions and multiple databases.

## Requirements

- Docker
- Docker Compose
- Make

## Usage

Place your project into the `projects` folder and run `make` to start the development environment.

An nginx server will be available at `http://<local-domain>` and the `projects` folder will be mounted into the server.

To set the virutal host, see the `./services/nginx/default/virtual-host-template.conf` file. You can copy the file and rename it to `<local-domain>.conf` and set the `server_name` to your desired local domain then add it to the `./services/nginx/sites-enabled` folder.

## Commands

- `make` - Start the development environment
- `make down` - Stop the development environment
- `make delete` - Remove all containers and volumes
- `make stop` - Stop all containers
- `make mysql8` - Access the MySQL 8 container
- `make mysql5` - Access the MySQL 5 container
- `make node` - Access the Node container (nvm to use multiple node versions)
- `make php82` - Access the PHP 8.2 container
- `make php81` - Access the PHP 8.1 container
- `make php80` - Access the PHP 8.0 container
- `make php74` - Access the PHP 7.4 container
