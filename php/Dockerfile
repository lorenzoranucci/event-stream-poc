FROM php:8.0.12-fpm

RUN apt update && apt install -y librdkafka-dev git zip unzip

RUN pecl install rdkafka && docker-php-ext-enable rdkafka

COPY --from=composer:2 /usr/bin/composer /usr/local/bin/composer

WORKDIR /srv

CMD ["php-fpm"]
