# go-irc2mail

深淵なる理由で IRC のメンションやキーワードをメールに通知したいので ~~heroku に~~ 鯖に bot を上げてキーワードを監視して通知出来るようにするものを作った

## 設定項目

環境変数で設定をする 項目は以下の通り

### MAILGUN_APIKEY


mailgun の API キー

### MAILGUN_DOMAIN

mailgun のドメイン

### IRC_NICKNAME

IRC サーバ接続に使用する nickname

### IRC_USERNAME

IRC サーバ接続に使用する username

### IRC_PASSWORD

IRC サーバ接続に使用する password

### IRC_ADDR

接続する IRCサーバのアドレス e.g. host[:port]

### IRC_USE_SSL

IRC サーバに SSL を使って接続するかどうか. 必要なら true や 1、不要なら false や 0 を指定する

### IRC_DEBUG

go-ircevent の debug 出力設定. 必要なら true や 1、不要なら false や 0 を指定する

### IRC_CHANNELS

join するチャンネル指定 スペース区切り

### IRC_MATCH_LIST

反応するキーワードの指定

### MAIL_FROM_NAME

メール送信時の From name

### MAIL_FROM_ADDR

メール送信時の From アドレス

### MAIL_TO_ADDR

メール送信時の送信先アドレス

### MAIL_SUBJECT

メール送信時の件名

### TARGET_NAME

監視する IRC ユーザ名(ログにしか使ってない)

### AWAY_NAME

監視する IRC ユーザの離席中のユーザ名

### APP_URL

heroku に上げてるときのアプリの URL

### DEBUG

デバッグ出力するかどうか. 必要なら true や 1、不要なら false や 0 を指定する

## 起動サンプル

```sh
#!/bin/sh

set -e

export MAILGUN_APIKEY="[mailgun apikey]"
export MAILGUN_DOMAIN="[mailgun domain]"

export IRC_NICKNAME="iku"
export IRC_USERNAME="iku"
export IRC_PASSWORD="******"
export IRC_ADDR="example.com[:port]"
export IRC_USE_SSL="true"
export IRC_DEBUG="true"
export IRC_CHANNELS="#channel1 #channel2"
export IRC_MATCH_LIST="all target"

export MAIL_FROM_NAME="target"
export MAIL_FROM_ADDR="noreply@example.com"
export MAIL_TO_ADDR="target@example.com"
export MAIL_SUBJECT="irc2mail"

export TARGET_NAME="target"
export AWAY_NAME="zz_target"
export APP_URL="http://localhost:19300"
export DEBUG="true"

export HOST=""
export PORT="19300"
export PING_MINUTE="15"

go run main.go
```

## Dockerfile サンプル

```Dockerfile
FROM ubuntu

MAINTAINER mix3

RUN apt-get update
RUN apt-get -y install wget jq curl unzip

WORKDIR /opt

RUN export ILLUSION_VERSION=`curl https://api.github.com/repos/mix3/go-irc2mail/releases | jq -r ".[0].tag_name"` && wget https://github.com/mix3/go-irc2mail/releases/download/${ILLUSION_VERSION}/go-irc2mail-${ILLUSION_VERSION}-linux-amd64.zip -O /opt/go-irc2mail.zip
RUN unzip go-irc2mail.zip
WORKDIR /opt/go-irc2mail

ENV MAILGUN_APIKEY="[mailgun apikey]"
ENV MAILGUN_DOMAIN="[mailgun domain]"

ENV IRC_NICKNAME="iku"
ENV IRC_USERNAME="iku"
ENV IRC_PASSWORD="******"
ENV IRC_ADDR="example.com[:port]"
ENV IRC_USE_SSL="true"
ENV IRC_DEBUG="true"
ENV IRC_CHANNELS="#channel1 #channel2"
ENV IRC_MATCH_LIST="all target"

ENV MAIL_FROM_NAME="target"
ENV MAIL_FROM_ADDR="noreply@example.com"
ENV MAIL_TO_ADDR="target@example.com"
ENV MAIL_SUBJECT="irc2mail"

ENV TARGET_NAME="target"
ENV AWAY_NAME="zz_target"
ENV APP_URL="http://localhost:19300"
ENV DEBUG="true"

ENV HOST=""
ENV PORT="19300"
ENV PING_MINUTE="1"

CMD ./go-irc2mail
```
