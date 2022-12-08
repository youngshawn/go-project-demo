#!/bin/sh

MYSQL_URL=127.0.0.1:3306

vault secrets enable database

vault write database/config/mysql-course \
    plugin_name=mysql-database-plugin \
    connection_url="{{username}}:{{password}}@tcp($MYSQL_URL)/" \
    allowed_roles="course-all" \
    username="root" \
    password="root"

vault write database/roles/course-all \
    db_name=mysql-course \
    creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT SELECT, INSERT, DELETE, UPDATE ON course.* TO '{{name}}'@'%';" \
    default_ttl="2m" \
    max_ttl="10m"


exit 0