# How to run

Start a docker container with our MySQL instance:
```
docker run --name mysql-test -e MYSQL_ROOT_PASSWORD=asd -e MYSQL_DATABASE=test --volume data:/var/lib/mysql -p 3306:3306 -d mysql:lts
```

Cleaning up is necessary.
Simply delete the container and the volume
```
docker rm --volumes mysql-test
```
