# [Choice] Go version (use -bullseye variants on local arm64/Apple Silicon): 1, 1.16, 1.17, 1-bullseye, 1.16-bullseye, 1.17-bullseye, 1-buster, 1.16-buster, 1.17-buster
ARG GO_VARIANT=1.18
ARG OPENCV_VARIANT=4.5.5

FROM ghcr.io/machinemapplatform/devcontainer-go${GO_VARIANT}-addlibkdslam-opencv-${OPENCV_VARIANT}:latest

# [Choice] Node.js version: none, lts/*, 16, 14, 12, 10
#ARG NODE_VERSION="none"
#RUN if [ "${NODE_VERSION}" != "none" ]; then su vscode -c "umask 0002 && . /usr/local/share/nvm/nvm.sh && nvm install ${NODE_VERSION} 2>&1"; fi

# [Optional] Uncomment this section to install additional OS packages.
# for docker-from-docker
RUN curl -sS https://dl.yarnpkg.com/debian/pubkey.gpg | sudo apt-key add - \
    echo "deb https://dl.yarnpkg.com/debian/ stable main" | sudo tee /etc/apt/sources.list.d/yarn.list \
    sudo apt update \
    sudo apt install yarn

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive \
    && apt-get -y install --no-install-recommends bash-completion pre-commit \
    cmake g++ wget unzip libgtk2.0-dev pkg-config glib-2.0 libcanberra-gtk* \
    libssl1.1 libkrb5-3 zlib1g libicu[0-9][0-9] gnome-keyring libsecret-1-0 desktop-file-utils x11-utils \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/*

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode
RUN go install github.com/golang/mock/mockgen@v1.6.0

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1

# install gitleaks
RUN curl -s https://api.github.com/repos/zricethezav/gitleaks/releases/latest |grep browser_download_url |cut -d '"' -f 4 |grep '\linux_x64' |wget -O gitleaks -i - \
    && tar --extract --file=gitleaks gitleaks \
    && chmod +x gitleaks && sudo mv gitleaks /usr/local/bin/
