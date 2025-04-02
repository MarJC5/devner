# Start from a base image with curl installed (Debian-based)
FROM debian:bullseye-slim

# Install dependencies
RUN apt-get update && apt-get install -y curl git bash build-essential libnss3-tools

# Install nvm
ENV NVM_DIR /usr/local/nvm

RUN mkdir -p $NVM_DIR

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash \
    && . "$NVM_DIR/nvm.sh"

ENV NODE_VERSION_18 18
ENV NODE_VERSION_20 20
ENV NODE_VERSION_22 22

# Use separate RUN commands to ensure that NVM is properly sourced each time
RUN . "$NVM_DIR/nvm.sh" && nvm install $NODE_VERSION_18 && nvm use $NODE_VERSION_18
RUN . "$NVM_DIR/nvm.sh" && nvm install $NODE_VERSION_20 && nvm use $NODE_VERSION_20
RUN . "$NVM_DIR/nvm.sh" && nvm install $NODE_VERSION_22 && nvm use $NODE_VERSION_22
RUN . "$NVM_DIR/nvm.sh" && nvm install --lts && nvm use --lts

# Optionally set a default Node version
RUN . "$NVM_DIR/nvm.sh" && nvm alias default $NODE_VERSION_20

# Verify installation
RUN . "$NVM_DIR/nvm.sh" && nvm use $NODE_VERSION_20 && node --version

# Install yarn
RUN . "$NVM_DIR/nvm.sh" && npm install -g yarn

# Install zsh & oh-my-zsh
RUN apt-get update && \
    apt-get install -y zsh && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

# Install Locale (for utf-8 support)
RUN apt-get update && \
    apt-get install -y \
        locales && \
    rm -r /var/lib/apt/lists/*

RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    locale-gen && \
    dpkg-reconfigure --frontend=noninteractive locales

ENV LC_ALL en_US.UTF-8 
ENV LANG en_US.UTF-8  
ENV LANGUAGE en_US:en

# Install postgres & mariadb clients
RUN apt-get update && \
    apt-get install -y mariadb-client postgresql-client && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Create directories for certificates
RUN mkdir -p /certs

# Install mkcert and generate certificates
RUN curl -JLO https://dl.filippo.io/mkcert/latest?for=linux/amd64 && \
    mv mkcert-v*-linux-amd64 /usr/local/bin/mkcert && \
    chmod +x /usr/local/bin/mkcert && \
    mkcert -install

# Install lsof
RUN apt-get update && \
    apt-get install -y lsof && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Copy your application source
WORKDIR /home/dev/apps/

# Keep the container running
CMD ["tail", "-f", "/dev/null"]