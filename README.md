Based on your provided `docker-compose` file, the `devner.sh` script, and the CLI menu script, here's an updated and comprehensive `README.md` for your Devner project:

```markdown
# Devner

## Description

Devner is a flexible development environment boilerplate designed to work with multiple PHP and Node.js versions, alongside various databases. It leverages Docker and Docker Compose to create an isolated and reproducible environment for development.

## Requirements

- Docker
- Docker Compose
- Make

## Setup and Usage

### 1. Clone the Repository

```sh
git clone https://github.com/MarJC5/devner.git
cd devner
```

### 2. Build and Start the Environment

```sh
make up
```

### 3. Access the Devner Menu

```sh
./devner.sh
```

### 4. Add alias to `.bashrc` or `.zshrc` (Optional)

To allow easy access to the Devner anywhere, you can add an alias to your shell configuration file:

```sh
./devner.sh alias

# Restart your terminal or run the following command:
source ~/.bashrc # or source ~/.zshrc

# Now you can access Devner from anywhere:
devner
```

## Devner Menu

The Devner menu allows you to manage various aspects of your development environment. Below are the available commands:

### General Commands

- **up**: Start the development environment.
- **down**: Stop and remove containers, networks, images, and volumes.
- **stop**: Stop all containers without removing them.
- **rebuild**: Rebuild the Docker containers.
- **delete**: Remove all containers and volumes.
- **nocache**: Remove all containers and volumes and remove cache.
- **reload**: Reload the frankenphp container to update Caddyfile.

### Database Commands

- **new-mysql**: Create a new MySQL database and user.
- **remove-mysql**: Remove a MySQL database and user.
- **new-postgres**: Create a new PostgreSQL database and user.
- **remove-postgres**: Remove a PostgreSQL database and user.

### Quick Access Commands

- **postgres**: Access the PostgreSQL container.
- **mysql**: Access the MySQL 8 container.
- **node**: Access the Node container.
- **frankenphp**: Access the FrankenPHP container.

### Project Commands

- **new**: Create a new Laravel or WordPress project.
- **remove**: Remove an existing project.

### Other Commands

- **ps**: Check if the devner container is running.
- **alias**: Add the devner alias to .bashrc or .zshrc.
- **code**: Open a project in VSCode.

### Caddyfile Commands

- **add-host**: Add a new host to the Caddyfile.
- **remove-host**: Remove an existing host from the Caddyfile.
- **list-hosts**: List all hosts in the Caddyfile.

### Help and Quit

- **credit**: Show information about the Devner tools and author.
- **help**: Show this help menu.
- **quit**: Exit the script.

## Example Usage

### Starting the Development Environment

To start the development environment, run:

```sh
devner up
```

### Creating a New Project

To create a new Laravel project with a MySQL database, run:

```sh
devner new laravel myproject mysql
```

### Removing a Project

To remove an existing project, run:

```sh
devner remove myproject mysql
```

## Docker Compose Configuration

The following services are defined in the `docker-compose.yml` file:

### Networks

- **default**: The default network for all services.

### Services

- **mailpit**: A local email testing service.
- **frankenphp**: A PHP server with Caddy and PHP-FPM.
- **node**: A Node.js development container.
- **mysql**: MySQL database server.
- **postgres**: PostgreSQL database server.
- **adminer**: Database management tool.

### Volumes

- **mysql_data**: Persistent storage for MySQL data.
- **postgres_data**: Persistent storage for PostgreSQL data.

## Author

[MarJC5](https://github.com/MarJC5)

Connect with me:

- **GitHub**: [yourusername](https://github.com/MarJC5)
- **LinkedIn**: [yourlinkedinprofile](https://linkedin.com/in/jean-christio-martin-385574111)
- **Twitter**: [@yourtwitterhandle](https://twitter.com/jeanchristio)

## License

This project is licensed under the MIT License.
