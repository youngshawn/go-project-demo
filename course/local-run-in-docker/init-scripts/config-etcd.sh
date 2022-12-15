#!/bin/sh

cat /remote-course.yaml | etcdctl --endpoints http://etcd-course:2379 put /config/prod/cloud/region/github.com/youngshawn/go-proect-demo/course/course.yaml


exit 0
