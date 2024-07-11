# Start from a base image with curl installed (Debian-based)
FROM debian:bullseye-slim

# Install dependencies
RUN apt-get update && apt-get install -y curl git bash build-essential

# Install nvm
ENV NVM_DIR /usr/local/nvm

RUN mkdir -p $NVM_DIR

RUN curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash \
    && . "$NVM_DIR/nvm.sh"

ENV NODE_VERSION_16 16
ENV NODE_VERSION_18 18
ENV NODE_VERSION_20 20

# Use separate RUN commands to ensure that NVM is properly sourced each time
RUN . "$NVM_DIR/nvm.sh" && nvm install $NODE_VERSION_16 && nvm use $NODE_VERSION_16
RUN . "$NVM_DIR/nvm.sh" && nvm install $NODE_VERSION_18 && nvm use $NODE_VERSION_18
RUN . "$NVM_DIR/nvm.sh" && nvm install $NODE_VERSION_20 && nvm use $NODE_VERSION_20
RUN . "$NVM_DIR/nvm.sh" && nvm install --lts && nvm use --lts

# Optionally set a default Node version
RUN . "$NVM_DIR/nvm.sh" && nvm alias default $NODE_VERSION_20

# Verify installation
RUN . "$NVM_DIR/nvm.sh" && nvm use $NODE_VERSION_20 && node --version

# Install yarn
RUN . "$NVM_DIR/nvm.sh" && npm install -g yarn

# Copy your application source
WORKDIR /home/dev/apps/

# Keep the container running
CMD ["tail", "-f", "/dev/null"]