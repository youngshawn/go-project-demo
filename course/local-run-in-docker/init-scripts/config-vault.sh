#!/bin/sh

# config database secret engine
vault secrets enable database
vault write database/config/mysql-course \
    plugin_name=mysql-database-plugin \
    connection_url="{{username}}:{{password}}@tcp(mysql-course:3306)/" \
    allowed_roles="course-all" \
    username="root" \
    password="root"
vault write database/roles/course-all \
    db_name=mysql-course \
    creation_statements="CREATE USER '{{name}}'@'%' IDENTIFIED BY '{{password}}'; GRANT CREATE, REFERENCES, SELECT, INSERT, DELETE, UPDATE ON course.* TO '{{name}}'@'%';" \
    default_ttl="2m" \
    max_ttl="10m"

# config transit secret engine
vault secrets enable transit
vault write -f transit/keys/course

# config approle auth engine
vault auth enable approle
vault policy write course -<<EOF
# Read-only permission on secrets stored at 'database/creds/course-all'
path "database/creds/course-all" {
  capabilities = [ "read" ]
}
# lease renew 相关权限
path "sys/leases/+/database/creds/course-all/*" {
  capabilities = [ "update" ]
}
# Transit相关权限
path "transit/encrypt/course" {
   capabilities = [ "update" ]
}
path "transit/decrypt/course" {
   capabilities = [ "update" ]
}
EOF
vault write auth/approle/role/course token_policies="course" \
    token_ttl=2m token_max_ttl=10m

# generate approle role-id and secret-id
vault read -field=role_id auth/approle/role/course/role-id >/vault/file/approle/course-roleid
vault write -f -field=secret_id auth/approle/role/course/secret-id >/vault/file/approle/course-secretid


exit 0
