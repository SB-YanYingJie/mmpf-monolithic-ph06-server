source ./env/devcontainer.env

if [[ $GITHUB_ACCESS_TOKEN != "" ]] ;
then
    echo $GITHUB_ACCESS_TOKEN | docker login ghcr.io -u $GITHUB_USERNAME --password-stdin
fi
