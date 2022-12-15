#!/bin/sh

cat /remote-course.yaml | etcdctl put /config/prod/cloud/region/github.com/youngshawn/go-proect-demo/course/course.yaml


exit 0
