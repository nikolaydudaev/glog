package parse

import (
	"regexp"
	"strings"
	"time"
)

type failbanRegex struct {
	Ip      *regexp.Regexp
	Flag    *regexp.Regexp
	Message *regexp.Regexp
	Action  *regexp.Regexp
}

var FailbanRegexInstance = failbanRegex{
	Ip:      regexp.MustCompile(`\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`),
	Flag:    regexp.MustCompile(`(fail2ban\.((actions)|(filter)) )`),
	Message: regexp.MustCompile(`((INFO)|(NOTICE))`),
	Action:  regexp.MustCompile(`((Found)|(Ban)|(Unban))`),
}

type failbanParser struct {
	FailbanEvent
}

var _ FailbanParser = (*failbanParser)(nil)

var FailbanParserInstance = failbanParser{}

// [parseFailbanTS] parses nginx timestamp string
// using layout "2006-01-02 15:04:05,000"
// and converts it to a string of the layout "2006-01-02 15:04:05"
func parseFailbanTS(s string) string {
	const layout = "2006-01-02 15:04:05,000"
	const target_layout = "2006-01-02 15:04:05"
	parsed, _ := time.Parse(layout, s)
	formatted := parsed.Format(target_layout)
	return formatted
}

func (f *failbanParser) FailbanParse(line string, out chan<- FailbanEvent) {
	failbanEvent := FailbanEvent{}
	ip := string(FailbanRegexInstance.Ip.Find([]byte(line)))
	failbanEvent.Ip = strings.TrimSpace(ip)
	date := strings.SplitN(line, " ", 3)
	failbanEvent.Ts = strings.Join([]string{date[0], date[1]}, " ")
	fl := string(FailbanRegexInstance.Flag.Find([]byte(line)))
	replacer := strings.NewReplacer("fail2ban.", "")
	fl = replacer.Replace(fl)
	fl = strings.TrimSpace(fl)
	failbanEvent.Flag = fl
	msg := string(FailbanRegexInstance.Message.Find([]byte(line)))
	msg = strings.TrimSpace(msg)
	failbanEvent.Message = msg
	action := string(FailbanRegexInstance.Action.Find([]byte(line)))
	action = strings.TrimSpace(action)
	failbanEvent.Action = action

	out <- failbanEvent
}
