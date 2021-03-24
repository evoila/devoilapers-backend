#!/bin/bash

echo "Enter CN (i.e. some.cert.de):"
read cn

openssl req -new -x509 -days 365 -nodes -text -out server.crt -keyout server.key -subj "/CN="$cn
chmod og-rwx server.key

openssl req -new -nodes -text -out root.csr -keyout root.key -subj "/CN="$cn
chmod og-rwx root.key

openssl x509 -req -in root.csr -text -days 3650 -extfile /etc/ssl/openssl.cnf -extensions v3_ca -signkey root.key -out root.crt

openssl req -new -nodes -text -out server.csr -keyout server.key -subj "/CN="$cn
chmod og-rwx server.key

openssl x509 -req -extfile <(printf "subjectAltName=DNS:$cn,DNS:www.$cn") -in server.csr -text -days 365 -CA root.crt -CAkey root.key -CAcreateserial -out server.crt