#!/usr/bin/env bash

cd /opt/mongodb
bin/mongod --dbpath=./data/ --fork --logpath=./data/fork.log

cd /opt/teaweb
bin/teaweb
