#!/bin/bash

echo "Loading docker images..."

if [ -e /home/mmpf/mmpf_modules/mmpf-service.tar ]; then
docker load < /home/mmpf/mmpf_modules/mmpf-service.tar &
pid1=$!
else
  echo "mmpf-service.tar not exists"
fi

if [ -e /home/mmpf/mmpf_modules/mmpf-other-services.tar ]; then
docker load < /home/mmpf/mmpf_modules/mmpf-other-services.tar &
pid2=$!
else
  echo "mmpf-other-services.tar not exists"
fi

wait $pid1 $pid2

mkdir -p /home/mmpf/mmpf_modules/mmpf-monolithic/lib
