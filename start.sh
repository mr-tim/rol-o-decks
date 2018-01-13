#!/bin/bash

FLASK_APP=server.py flask run > /dev/null 2>&1 &
echo $! > flask.pid

pushd ui
node_modules/.bin/webpack-dev-server --inline --progress --config build/webpack.dev.conf.js > /dev/null 2>&1 &
NODE_PID=$!
popd
echo $NODE_PID>node.pid

python index_pptx.py > /dev/null 2>&1 &
echo $!>index.pid