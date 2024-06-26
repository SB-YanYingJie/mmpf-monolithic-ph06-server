#!/bin/bash

ENV_DIR=./devices
ENV_LIST_FILE=./devices/env_list.txt

if [ "$#" == 0 ]; then
  echo "Specify the env file name as argument"
  exit 1
fi
# COMPOSE_SERVICE_FILE=$1
PARAM_FILE=$1
echo "param file: ""$PARAM_FILE"

# read envfile list
ENV_FILE_PATHS=()
IFS=$'\n'
TMP_FILE=/tmp/__env

echo "--- env files ---"
while read -r ENV_FILE || [[ -n "${ENV_FILE}" ]];do
  # コメント行除外
  ENV_FILE_PATH=$ENV_DIR/"$ENV_FILE"
  echo "$ENV_FILE" | grep -v '^#.*' > /dev/null
  if [ $? -eq 0 ];then
    if [ -f "$ENV_FILE_PATH" ];then
      echo "$ENV_FILE_PATH"
      ENV_FILE_PATHS+=( "$ENV_FILE_PATH" )
    else
      echo "$ENV_FILE_PATH" "is not found"
      exit 1
    fi
  fi
done < ${ENV_LIST_FILE}

echo "--- services ---"
for ENV_FILE_PATH in "${ENV_FILE_PATHS[@]}"; do
  if [ -f "$ENV_FILE_PATH" ]; then
    echo "user env file: ""$ENV_FILE_PATH"

    echo "`cat "$ENV_FILE_PATH"`" > $TMP_FILE
    echo "`cat "$PARAM_FILE"`" >> $TMP_FILE
    env ENV_PATH=$TMP_FILE docker-compose -f docker-compose.service.mmpf.yml -p "${ENV_FILE_PATH##*.}" --env-file $TMP_FILE up --force-recreate -d
  fi
done
rm -f /tmp/__env
