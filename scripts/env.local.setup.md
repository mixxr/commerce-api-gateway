## download and create a local instance of MariaDB docker:
[MariaDB in a Docker, Complete Guide](https://mariadb.com/kb/en/installing-and-using-mariadb-via-docker/)

if you had Docker running on your local machine, this command download, create and instance the latest mariadb docker image:

`docker container run \
        --name dbmarialocal \
        -e MYSQL_ROOT_PASSWORD=secr3tZ \
        -e MYSQL_USER=golang \
        -e MYSQL_PASSWORD=secr3tZ \
        -e MYSQL_DATABASE=dcgw \
        -p 3306:3306 \
        -d mariadb/server`

### to accept connections:
/etc/mysql/my.cnf -> 

 ## install Golang on Mac via brew

`brew update; brew upgrade; brew install go --cross-compile-common`

`go get github.com/go-sql-driver/mysql`

