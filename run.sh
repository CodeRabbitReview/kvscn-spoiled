#!/bin/bash

docker run -v $(pwd)/persistence:/root/persistence -p 8080:8080 \
--add-host=host.docker.internal:host-gateway --name storage_server --rm miprokop/storage_server & \
bash sender.sh test.json ok