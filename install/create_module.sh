#!/bin/bash

# mkdir tmp dir
if [ ! -d "/tmp/mmpf" ]; then
    mkdir /tmp/mmpf
else
    rm -rf /tmp/mmpf/*
fi

TAG=${1:-latest}

echo "docker pull..."
# image mmpf
docker pull ghcr.io/machinemapplatform/mmpf-monolithic:$TAG

# image mmpf-server-services
docker pull redis:6.2.1

echo "Make image in tar format..."
docker save -o /tmp/mmpf/mmpf-service.tar ghcr.io/machinemapplatform/mmpf-monolithic:$TAG
docker save -o /tmp/mmpf/mmpf-other-services.tar redis

echo "Clone mmpf-monolithic repositiory..."
cd /tmp/mmpf && git clone https://github.com/machinemapplatform/mmpf-monolithic.git


