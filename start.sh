#! /bin/bash

docker run --name mysql-test -e MYSQL_ROOT_PASSWORD=asd -e MYSQL_DATABASE=test --volume data:/var/lib/mysql --volume $PWD/conf:/etc/mysql/conf.d -p 3306:3306 -d mysql:lts
