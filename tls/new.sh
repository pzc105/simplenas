#!/bin/bash

if [[ $# > 0 && $1 = '-ca' ]] 
then
  echo "gen ca"
  openssl genrsa -out ca.key 2048
  openssl req -new -x509 -days 3650 -key ca.key -subj "/C=CN/ST=GD/L=SZ/O=PNas, Inc./CN=PNas CA" -out ca.crt
elif [ $# -eq 2 ]
then
  openssl req -newkey rsa:2048 -nodes -keyout $1.key -subj "/C=CN/ST=GD/L=SZ/O=PNas, Inc./CN=$2" -out $1.csr
  openssl x509 -req -extfile <(printf "subjectAltName=DNS:$2") -days 3650 -in $1.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out $1.crt
fi
