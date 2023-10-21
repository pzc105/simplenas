#!/bin/bash

crt_path="./http.crt"
key_path="./http.key"
while getopts "c:k:" opt
do
  case $opt in
    c)
      crt_path=$OPTARG
      ;;
    k)
      key_path=$OPTARG
      ;;
    ?)
      echo "getopts param error"
      exit 1;;
  esac
done
echo crt_path: ${crt_path}
echo key_path: ${key_path}
PORT=3000 HTTPS=true SSL_CRT_FILE=${crt_path} SSL_KEY_FILE=${key_path} (nohup npm start &)