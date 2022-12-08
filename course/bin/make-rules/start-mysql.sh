#!/bin/sh

docker pull mysql/mysql-server:5.7

mkdir -p ~/data/mysql-data

docker run --name mysql-course \
        --publish 3306:3306 \
        --volume ~/data/mysql-data:/var/lib/mysql \
        --env MYSQL_ROOT_PASSWORD=root \
        --env MYSQL_ROOT_HOST=% \
        --env MYSQL_DATABASE=course \
        --env MYSQL_USER=course_user \
        --env MYSQL_USER_HOST=% \
        --env MYSQL_PASSWORD=c_u_passw0rd \
        --detach mysql/mysql-server:5.7

if [ $? -ne 0 ] then
    echo "[error] docker run failed, exit."
    exit 1
fi

exit 0