#! /bin/bash

docker stop mysql-test
docker rm mysql-test
docker volume rm data
