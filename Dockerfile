FROM ubuntu:18.04

ENV GO_VERSION 1.11.2
ENV OS linux
ENV ARCH amd64
ENV GOROOT /usr/local/go
ENV GOPATH /go
ENV PORT 8008
ENV PATH /usr/local/go/bin:/go/bin:$PATH

# 必要なライブラリの導入
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y curl git make sudo
RUN curl -O https://dl.google.com/go/go$GO_VERSION.$OS-$ARCH.tar.gz && \
    tar -C /usr/local -xzf go$GO_VERSION.$OS-$ARCH.tar.gz && \
    go version && \
    rm go$GO_VERSION.$OS-$ARCH.tar.gz

# ユーザの作成
RUN yes | adduser --disabled-password echouser

# goのモジュール群のディレクトリを作成
RUN mkdir -p /go/bin /go/src/github.com/LAKuEN/echoproject && chmod -R 777 /go/bin
WORKDIR /go/src/github.com/LAKuEN/echoproject
COPY . .


# gocvの導入
# FIXME Herokuにデプロイする為にはOpenCVのインストールパスを変える必要があるはず
#       その場合はgocvで参照するOpenCVのパスも変えねばならない？もしくはパスが通っていれば問題ない？
RUN go get -u -d gocv.io/x/gocv
# gocvの導入
RUN cd $GOPATH/src/gocv.io/x/gocv && \
# FIXME make -j8とかにした方が良いかも。めちゃくちゃ遅い
# FIXME ./cmd/version/main.goが実行できないエラーが出る
#       どのディレクトリに居るかを出力してやればいいかも
make -j8 sudo_install
# アプリのインストール
RUN go install github.com/LAKuEN/echoproject

# サーバ起動
USER echouser
CMD echoproject
