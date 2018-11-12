FROM ubuntu:18.04

ENV GO_VERSION 1.11.2
ENV OS linux
ENV ARCH amd64
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PORT 8008
ENV PATH /usr/local/go/bin:/go/bin:$PATH

# 各種ライブラリの導入
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y curl git make sudo
RUN curl -O https://dl.google.com/go/go$GO_VERSION.$OS-$ARCH.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.$OS-$ARCH.tar.gz && \
    go version && \
    rm go$GO_VERSION.$OS-$ARCH.tar.gz

RUN yes | adduser --disabled-password cvuser

# goのモジュール群のディレクトリを作成
RUN mkdir -p /go/bin /go/src/github.com/LAKuEN/cut-and-draw-stickynote-api && chmod -R 777 /go/bin
WORKDIR /go/src/github.com/LAKuEN/cut-and-draw-stickynote-api
COPY . .

# gocvの導入
RUN go get -u -d gocv.io/x/gocv
RUN cd $GOPATH/src/gocv.io/x/gocv && \
make -j8 sudo_install
# アプリのインストール
RUN go install github.com/LAKuEN/cut-and-draw-stickynote-api

# サーバ起動
USER cvuser
CMD cut-and-draw-stickynote-api
