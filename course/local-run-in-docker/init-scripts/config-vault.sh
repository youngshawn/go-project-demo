#!/bin/sh

# apk add
apk add jq

# init and unseal
if vault status | grep -q "Initialized.*false" ; then
    # vault init
    vault operator init -format json >/vault/file/vault-keys/keys.json
fi
if vault status | grep -q "Sealed.*true" ; then
    # vault unseal
    vault operator unseal $(cat /vault/file/vault-keys/keys.json | jq -r .unseal_keys_b64[1])
    vault operator unseal $(cat /vault/file/vault-keys/keys.json | jq -r .unseal_keys_b64[2])
    vault operator unseal $(cat /vault/file/vault-keys/keys.json | jq -r .unseal_keys_b64[3])
fi

# set root_token env
export VAULT_TOKEN=$(cat /vault/file/vault-keys/keys.json | jq -r .root_token)

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
# lease renew
path "sys/leases/+/database/creds/course-all/*" {
  capabilities = [ "update" ]
}
# Transit
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
