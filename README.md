# go-irc2mail

深淵なる理由で IRC のメンションやキーワードをメールに通知したいので VPS に bot を上げてキーワードを監視して通知出来るようにするものを作った

## Usage

### Gmail API

メール送信 に Gmail API を使用するので事前に OAuth2 の ClientID, ClientSecret, RefreshToken を取得しておく

* [Google Developers Console](https://console.developers.google.com/)へ行く
 * 適当なプロジェクトが無ければ作成する
* Gmail API を有効にする
 * API > Gmail API > APIを有効にする
* 認証情報を追加する
 * 認証情報 > 認証情報を追加 > OAuth 2.0 クライアント ID > その他

この時点で ClientID, と ClientSecret が取得出来ている。

RefreshToken は取得するためのコマンドを用意したので、以下のようにして指示に従って進めると取得可能

```
go get github.com/mix3/go-irc2mail/cmd/rtkn
rtkn --client_id ****** --client_secret ******
.
.
.
```

### irc2mail

設定が必要なので config.yaml.sample をコピって OAuth2 の設定やら IRC 接続の設定やら通知の設定やらをする

```yaml
# gmail api でメール送信するための設定
oauth2:
  client_id:     "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX.apps.googleusercontent.com"
  client_secret: "XXXXXXXXXXXXXXXXXXXXXXXX"
  refresh_token: "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

# bot の irc ログイン設定
ircevent:
  nickname: "botname"
  username: "botname"
  password: "******"
  host:     "irc.host"
  port:     6667
  use_ssl:  false
  debug:    false

# 反応するチャンネル、キーワードの通知設定
notice:
  "botname":
    "#channel":
      keywords:
        - "all"
        - "botname"
#       - .
#       - .
#       - .
      to:             "notice@example.com"
      subject_prefix: "irc2mail"
      check_away:     false
      away_name:      "zz_botname"

# アプリログ出力設定
debug_log: false
```

設定したら以下のようにして起動する

```
go get github.com/mix3/go-irc2mail/cmd/irc2mail
irc2mail -c config.yaml # default で config.yaml を見るので -c ****** は省略可
```

## LICENSE

MIT
