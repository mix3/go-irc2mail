package irc2mail

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChannels(t *testing.T) {
	b := Bot{
		Config: Config{
			Notice: ConfigNoticeMap{
				"target_1": {
					"#channel_1": ConfigNotice{},
					"#channel_2": ConfigNotice{},
				},
				"target_2": {
					"#channel_2": ConfigNotice{},
					"#channel_3": ConfigNotice{},
				},
				"target_3": {
					"#channel_3": ConfigNotice{},
					"#channel_4": ConfigNotice{},
				},
			},
		},
	}
	for _, expect := range []string{
		"#channel_1",
		"#channel_2",
		"#channel_3",
		"#channel_4",
	} {
		assert.Contains(t, b.Channels(), expect)
	}
}
