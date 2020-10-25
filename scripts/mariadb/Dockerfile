FROM mariadb:latest as builder

LABEL author="ms.co.it@gmail.com"
# if sample data are needed
ARG sampledata="true" 
RUN echo "sampledata=$sampledata" > /tmp/setup.txt
# That file does the DB initialization but also runs mysql daemon, by removing the last line it will only init
# RUN ["sed", "-i", "s/exec \"$@\"/echo \"not running $@\"/", "/usr/local/bin/docker-entrypoint.sh"]

# needed for intialization
ENV MYSQL_ROOT_PASSWORD=secr3tZ
ENV MYSQL_USER=golang
ENV MYSQL_PASSWORD=secr3tZuser
ENV MYSQL_DATABASE=dcgw

EXPOSE 3306/tcp

COPY init.sh /docker-entrypoint-initdb.d/
COPY setup.sql /docker-entrypoint-initdb.d/
COPY sample.sql /tmp/
RUN if [ "$sampledata" = "true" ]; then echo "source /tmp/sample.sql" >> /docker-entrypoint-initdb.d/setup.sql; fi

# Need to change the datadir to something else that /var/lib/mysql because the parent docker file defines it as a volume.
# https://docs.docker.com/engine/reference/builder/#volume :
#       Changing the volume from within the Dockerfile: If any build steps change the data within the volume after
#       it has been declared, those changes will be discarded.
#RUN ["/usr/local/bin/docker-entrypoint.sh", "mysqld", "--datadir", "/initialized-db", "--aria-log-dir-path", "/initialized-db"]

#FROM mariadb:latest

# COPY --from=builder /initialized-db /var/lib/mysql