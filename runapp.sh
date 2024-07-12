#!/bin/bash

docker compose run --rm -v $PWD:/opt/srv/api -p 8080:8080 app /bin/ash

