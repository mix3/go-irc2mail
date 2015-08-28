package irc2mail

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ConfigOAuth2 struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RefreshToken string `yaml:"refresh_token"`
}

type ConfigIRCEvent struct {
	NickName string `yaml:"nickname"`
	UserName string `yaml:"username"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UseSSL   bool   `yaml:"use_ssl"`
	Debug    bool   `yaml:"debug"`
}

type ConfigNotice struct {
	Keywords      []string `yaml:"keywords"`
	To            string   `yaml:"to"`
	SubjectPrefix string   `yaml:"subject_prefix"`
	CheckAway     bool     `yaml:"check_away"`
	AwayName      string   `yaml:"away_name"`
}

type ConfigNoticeMap map[string]map[string]ConfigNotice

type Config struct {
	OAuth2   ConfigOAuth2    `yaml:"oauth2"`
	IRCEvent ConfigIRCEvent  `yaml:"ircevent"`
	Notice   ConfigNoticeMap `yaml:"notice"`
	DebugLog bool            `yaml:"debug_log"`
}

func LoadConfig(path string) (Config, error) {
	var c Config
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return c, err
	}
	if err := yaml.Unmarshal(data, &c); err != nil {
		return c, err
	}
	return c, nil
}
