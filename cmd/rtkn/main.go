package main

import (
	"fmt"
	"log"

	"github.com/jessevdk/go-flags"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
)

func main() {
	var opts struct {
		ClientID     string `long:"client_id"     required:"true"`
		ClientSecret string `long:"client_secret" required:"true"`
	}
	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatalln(err)
		return
	}

	config := &oauth2.Config{
		ClientID:     opts.ClientID,
		ClientSecret: opts.ClientSecret,
		RedirectURL:  "urn:ietf:wg:oauth:2.0:oob",
		Scopes:       []string{"https://www.googleapis.com/auth/gmail.send"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://accounts.google.com/o/oauth2/auth",
			TokenURL: "https://accounts.google.com/o/oauth2/token",
		},
	}

	var code string
	url := config.AuthCodeURL("")
	fmt.Println("ブラウザで以下のURLにアクセスし、認証してCodeを取得してください。")
	fmt.Println(url)
	fmt.Println("取得したCodeを入力してください")
	fmt.Scanf("%s\n", &code)
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		log.Fatalln("Exchange: ", err)
	}
	fmt.Println("RefreshToken: ", token.RefreshToken)
}
