## 1. Create `openssl.cnf`

```toml
[req]
distinguished_name = req_distinguished_name
x509_extensions = v3_req
prompt = no

[req_distinguished_name]
CN = localhost

[v3_req]
keyUsage = keyEncipherment, dataEncipherment
extendedKeyUsage = serverAuth
subjectAltName = @alt_names

[alt_names]
DNS.1 = localhost
```

## 2. Create certificate

```shell
openssl req -x509 -sha256 -nodes -days 99999 -newkey rsa:2048 -keyout localhost.key -out localhost.crt -config openssl.cnf
```
