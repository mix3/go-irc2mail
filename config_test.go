package irc2mail

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TempFile() (string, string, error) {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", "", err
	}

	file, err := ioutil.TempFile(dir, "")
	if err != nil {
		return "", "", err
	}

	return file.Name(), dir, nil
}

func TestLoadConfig(t *testing.T) {
	tmpfile, tmpdir, err := TempFile()
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	err = ioutil.WriteFile(tmpfile, []byte(`
# gmail api でメール送信するための設定
oauth2:
  client_id:     "[CLIENT_ID]"
  client_secret: "[CLIENT_SECRET]"
  refresh_token: "******"

# 自分の nick 監視用 bot の設定
ircevent:
  nickname: "bot"
  username: "bot"
  password: "******"
  host:     "irc.host.name"
  port:     6667
  use_ssl:  false
  debug:    false

# 反応するチャンネル、キーワードの設定
notice:
  "target_1":
    "#channel_1":
      keywords:
        - "all"
        - "target_1"
      to:             "target_1@example.com"
      subject_prefix: "subject_prefix"
      check_away:     false
      away_name:      ""
    "#channel_2":
      keywords:
        - "all"
        - "target_1"
      to:             "target_1@example.com"
      subject_prefix: "subject_prefix"
  "target_2":
    "#channel_1":
      keywords:
        - "all"
        - "target_2"
      to:             "target_2@example.com"
      subject_prefix: "subject_prefix"
      check_away:     true
      away_name:      "zz_target_2"
`), 0644)
	if err != nil {
		t.Fatal(err)
	}

	expect := Config{
		OAuth2: ConfigOAuth2{
			ClientID:     "[CLIENT_ID]",
			ClientSecret: "[CLIENT_SECRET]",
			RefreshToken: "******",
		},
		IRCEvent: ConfigIRCEvent{
			NickName: "bot",
			UserName: "bot",
			Password: "******",
			Host:     "irc.host.name",
			Port:     6667,
			UseSSL:   false,
			Debug:    false,
		},
		Notice: ConfigNoticeMap{
			"target_1": {
				"#channel_1": ConfigNotice{
					Keywords: []string{
						"all",
						"target_1",
					},
					To:            "target_1@example.com",
					SubjectPrefix: "subject_prefix",
					CheckAway:     false,
					AwayName:      "",
				},
				"#channel_2": ConfigNotice{
					Keywords: []string{
						"all",
						"target_1",
					},
					To:            "target_1@example.com",
					SubjectPrefix: "subject_prefix",
					CheckAway:     false,
					AwayName:      "",
				},
			},
			"target_2": {
				"#channel_1": ConfigNotice{
					Keywords: []string{
						"all",
						"target_2",
					},
					To:            "target_2@example.com",
					SubjectPrefix: "subject_prefix",
					CheckAway:     true,
					AwayName:      "zz_target_2",
				},
			},
		},
	}

	got, err := LoadConfig(tmpfile)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, got, expect)
}
