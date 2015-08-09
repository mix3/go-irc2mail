# go-irc2mail

深淵なる理由で IRC のメンションやキーワードをメールに通知したいので heroku に bot を上げてキーワードを監視して通知出来るようにするものを作った

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
