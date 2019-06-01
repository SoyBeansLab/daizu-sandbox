FROM golang:1.12.5-alpine3.9

MAINTAINER Taichi Uchihara <nil.uchihara@gmail.com>

RUN \
  apk add git && \
  git clone https://github.com/SoyBeansLab/daizu-sandbox /daizu-sandbox && \
  cd /daizu-sandbox && git checkout develop && git pull && \
  go get -u 

WORKDIR /daizu-sandbox

ENTRYPOINT ["sh"]
