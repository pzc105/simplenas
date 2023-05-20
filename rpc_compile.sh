#!/bin/bash

cd ./src/frontend
python3 rpc_compile.py

cd -
cd ./src/server
python3 rpc_compile.py

cd -
cd ./src/bt
python3 rpc_compile.py