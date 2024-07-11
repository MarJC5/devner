FROM dunglas/frankenphp

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
    mysqli \
    pgsql \
    gd \
    intl \
    zip \
    opcache

RUN apt-get update && \
    apt-get install -y \
    postfix \
    mailutils \
    libnss3-tools \
    git \
    openssh-client && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install MailPit (SMTP testing tool)
RUN curl -sL https://raw.githubusercontent.com/axllent/mailpit/develop/install.sh | bash

# Install wp-cli
RUN curl -O https://raw.githubusercontent.com/wp-cli/builds/gh-pages/phar/wp-cli.phar && \
    chmod +x wp-cli.phar && \
    mv wp-cli.phar /usr/local/bin/wp

# Install Composer
RUN curl -sS https://getcomposer.org/installer | php && \
    mv composer.phar /usr/local/bin/composer


# Set the working directory
WORKDIR /var/www/html
