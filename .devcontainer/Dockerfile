FROM golang:1.25.4

RUN apt-get update && apt-get install -y \
    curl git bash docker.io \
    && rm -rf /var/lib/apt/lists/*

# Install NVM
ENV NVM_DIR=/usr/local/nvm
ENV NODE_VERSION=lts
RUN mkdir -p $NVM_DIR \
    && curl -o- https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.5/install.sh | bash

ENV NODE_PATH=$NVM_DIR/versions/node/$NODE_VERSION/lib/node_modules
ENV PATH=$NVM_DIR/versions/node/$NODE_VERSION/bin:$PATH
