#!/bin/bash

kill -9 `cat index.pid`
rm index.pid
kill -9 `cat node.pid`
rm node.pid
kill -9 `cat flask.pid`
rm flask.pid