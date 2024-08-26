package configs

import "strings"

type skips []string

var (
	HTTP                       = false
	BaseUrl                    = ""
	Project                    = ""
	ChangelogHeaderLines       = 2
	MR                         = false
	SkipMsgs             skips = make([]string, 0)
)

func ParseSkipMsg(msg string) error {
	msgs := strings.Split(msg, ",")
	for _, m := range msgs {
		SkipMsgs = append(SkipMsgs, m)
	}
	return nil
}

func (s *skips) ShouldSkip(msg string) bool {
	for _, m := range *s {
		if strings.TrimSpace(m) == strings.TrimSpace(msg) {
			return true
		}
	}
	return false
}
