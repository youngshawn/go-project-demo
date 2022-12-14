# build from alpine
FROM alpine:latest

# make workdir
RUN mkdir /app
WORKDIR /app/

# copy binary
COPY dist/course-linux /app/course

# volume for conf and approle
VOLUME [ "/app/conf", "/app/role" ]

# expose port 8080
EXPOSE 8080

# entrypoint and cmd
ENTRYPOINT [ "/app/course", "--listen", ":8080" ]
CMD [ "--config", "/app/conf/course.yaml", "--vault.auth.roleid-file-path", "/app/role/course-roleid", "--vault.auth.secretid-file-path", "/app/role/course-secretid" ]