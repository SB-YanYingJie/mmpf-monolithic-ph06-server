#!/bin/bash

DOCKER_NETWORK_NAME=monolithic_network


#shutdown devcontainer-redis
if [ -n "$(docker ps -aq --filter "name=mmpf-monolithic_devcontainer-redis-1")" ]; then
    echo "remove devcontainer-redis container"
    docker rm -f $(docker ps --filter "name=mmpf-monolithic_devcontainer-redis-1")
else
    echo "no devcontainer-redis container"
fi

# create network
NW_CNT=$(docker network ls -q --filter name=$DOCKER_NETWORK_NAME | wc -l)
if [ "$NW_CNT" -eq 0 ]; then
  echo "create network:$DOCKER_NETWORK_NAME"
  docker network create $DOCKER_NETWORK_NAME
fi
wait

# run redis
docker-compose -f docker-compose.redis.yml up -d
# ↓ EC2でサーバー起動した際、これが無いとredisに画像を書き込めないことがある
docker exec mmpf-monolithic_redis_1 chown -R redis:redis /etc

# run mmpf
./scripts/start_mmpf.sh ./devices/common.env
