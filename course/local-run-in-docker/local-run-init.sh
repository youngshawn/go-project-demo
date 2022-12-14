#!/bin/sh

vault read -field=role_id auth/approle/role/course/role-id >local-run-in-docker/approle/course-roleid
vault write -f -field=secret_id auth/approle/role/course/secret-id >local-run-in-docker/approle/course-secretid

cat remote-course.yaml | etcdctl put /config/prod/cloud/region/github.com/youngshawn/go-proect-demo/course/course.yaml

exit 0