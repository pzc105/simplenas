#!/bin/bash

while ! mysqladmin ping -h"127.0.0.1" --silent; do
  sleep 1
done