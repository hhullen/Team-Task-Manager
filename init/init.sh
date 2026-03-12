#!/bin/bash
set -e

# Читаем секреты
MIGRATOR_USER=$(cat /run/secrets/db_migrator_user)
MIGRATOR_PASS=$(cat /run/secrets/db_migrator_password)
APP_USER=$(cat /run/secrets/db_app_user)
APP_PASS=$(cat /run/secrets/db_app_password)
DB_NAME=$(cat /run/secrets/db_name)

TEMPLATE="/sql_init/init_roles.sql"

sed -e "s|{{MIGRATOR_USER}}|$MIGRATOR_USER|g" \
    -e "s|{{MIGRATOR_PASS}}|$MIGRATOR_PASS|g" \
    -e "s|{{APP_USER}}|$APP_USER|g" \
    -e "s|{{APP_PASS}}|$APP_PASS|g" \
    -e "s|{{DB_NAME}}|$DB_NAME|g" \
    "$TEMPLATE" | mysql -u root -p"$MYSQL_ROOT_PASSWORD"