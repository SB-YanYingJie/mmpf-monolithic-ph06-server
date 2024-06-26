#!/bin/bash

source /app/env/devcontainer.env

sudo chown -R vscode:vscode /app
sudo chown -R vscode:vscode /go

safedir=`git config safe.directory`
if [[ $safedir != "/app" ]]; then
    git config --add safe.directory /app
fi

url=`git config --get remote.origin.url`
if [[ $url = git@* ]]; then
    git config url.ssh://git@github.com/.insteadOf https://github.com/
elif [[ $GITHUB_USERNAME != "" ]]; then
    # https://go.dev/doc/faq#git_https:~:text=how%20to%20proceed.-,Why%20does%20%22go%20get%22%20use%20HTTPS%20when%20cloning%20a%20repository%3F,-Companies%20often%20permit
    touch ~/.netrc
    echo "machine github.com login $GITHUB_USERNAME password $GITHUB_ACCESS_TOKEN" > ~/.netrc
    cp ~/.netrc /app/mock-server/.netrc # for building mock-server
    cp ~/.netrc /app/.netrc # for building app
fi

aws configure set region ap-northeast-1
aws configure set aws_access_key_id --profile dev $AWS_ACCESS_KEY_ID_DEV
aws configure set aws_secret_access_key --profile dev $AWS_SECRET_ACCESS_KEY_DEV
pre-commit install

sudo chown vscode:vscode /tmp/*.so /tmp/*.kdbow /tmp/*.kdmp
mkdir -p /app/lib
sudo mv -f /tmp/libKdSlam.so /usr/local/lib/
chmod -R +x /app/lib/
mv -f /tmp/*.kdmp /app/lib/
mv -f /tmp/*.kdbow /app/lib/

go mod tidy
