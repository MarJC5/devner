FROM dunglas/frankenphp

# Set the timezone
ENV TZ=UTC
RUN echo $TZ > /etc/timezone && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata

# Install Node.js and NVM
ENV NVM_DIR /usr/local/nvm
RUN mkdir -p $NVM_DIR && \
    curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash

ENV NODE_VERSION_18 18
ENV NODE_VERSION_20 20
ENV NODE_VERSION_22 22

# Install Node.js versions
RUN . "$NVM_DIR/nvm.sh" && \
    nvm install $NODE_VERSION_18 && \
    nvm install $NODE_VERSION_20 && \
    nvm install $NODE_VERSION_22 && \
    nvm install --lts && \
    nvm alias default $NODE_VERSION_22 && \
    nvm use $NODE_VERSION_22 && \
    npm install -g yarn && \
    npm install -g pnpm
    
# Set the timezone
ENV TZ=UTC
RUN echo $TZ > /etc/timezone && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    dpkg-reconfigure -f noninteractive tzdata


# Install Postfix and PHP extensions
RUN install-php-extensions \
    pdo \
    pdo_mysql \
    pdo_pgsql \
    redis \
    mysqli \
    pgsql \
    gd \
    intl \
    zip \
    opcache \
    pcntl \
    exif \
    bcmath

RUN apt-get update && \
    apt-get install -y \
    mariadb-client \
    postgresql-client \
    postfix \
    mailutils \
    libnss3-tools \
    git \
    openssh-client && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install wp-cli
RUN curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar && \
    chmod +x wp-cli.phar && \
    mv wp-cli.phar /usr/local/bin/wp

# Install Composer
RUN curl -sS https://getcomposer.org/installer | php && \
    mv composer.phar /usr/local/bin/composer

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

# Install MailPit (SMTP testing tool)
RUN curl -sL https://raw.githubusercontent.com/axllent/mailpit/develop/install.sh | bash

RUN export PATH="$PATH:$HOME/.composer/vendor/bin" >> ~/.bashrc

# Set the working directory
WORKDIR /var/www/html