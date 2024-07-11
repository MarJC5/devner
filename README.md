# Devner

## Description

Devner is a flexible development environment boilerplate designed to work with multiple PHP and Node.js versions, alongside various databases. It leverages Docker and Docker Compose to create an isolated and reproducible environment for development.

## Requirements

- Docker
- Docker Compose
- Make
- Visual Studio Code (for the `code` command)

## Setup and Usage

1. Clone this repository to your local machine.
2. Place your projects into the `projects` folder.
3. Run `./devner.sh up` to start the development environment. This command boots up the necessary Docker containers for your development needs.

An Nginx server will be made available at `http://<local-domain>`, and the `projects` folder will be mounted into the server for easy file access and live edits.

### Setting Up Virtual Hosts

To set up a virtual host for a project:

1. Navigate to the `./services/nginx/default/` directory.
2. Use the `virtual-host-template.conf` file as a template.
3. Copy and rename the template to `<your-local-domain>.conf`.
4. Modify the `server_name` directive to match your desired local domain.
5. Place the modified file in the `./services/nginx/sites-enabled` folder for it to take effect.

Don't forget to add your local domain to your `/etc/hosts` file pointing to `127.0.0.1`.

## Commands

Run these commands from the root of the cloned repository:

- `./devner.sh up` - Start the development environment.
- `./devner.sh down` - Stop and remove containers, networks, images, and volumes.
- `./devner.sh stop` - Stop all containers without removing them.
- `./devner.sh rebuild` - Rebuild the Docker containers.
- `./devner.sh delete` - Remove all containers and volumes.
- Database access commands:
  - `./devner.sh mysql8` - Access the MySQL 8 container.
  - `./devner.sh mysql5` - Access the MySQL 5 container.
- Node and PHP version-specific commands:
  - `./devner.sh node` - Access the Node container (uses NVM for multiple Node.js versions).
  - `./devner.sh php82` - Access the PHP 8.2 container.
  - `./devner.sh php81` - Access the PHP 8.1 container.
  - `./devner.sh php8` - Access the PHP 8.0 container.
  - `./devner.sh php74` - Access the PHP 7.4 container.
- Additional commands:
  - `./devner.sh host` - Add a new host entry to `/etc/hosts` and create a new virtual host.
  - `./devner.sh reload` - Reload Nginx configuration without restarting the container.
  - `./devner.sh code` - Open a project in Visual Studio Code by listing available projects in the `projects` folder. Select a project by entering the corresponding number when prompted.
  - `./devner.sh add-alias` - Add an alias to the `./devner.sh` script for easier access.
  - `./devner.sh ps` - Check the status of the Docker containers.
- New project setup:
  - `./devner.sh new <laravel/wp> <project-name>` - Create a new Laravel or WordPress project in the `projects` folder.
  