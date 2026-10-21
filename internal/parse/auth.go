package parse

import (
	"regexp"
	"strings"
	"time"
)

type authRegex struct {
	Ip     *regexp.Regexp
	Port   *regexp.Regexp
	Reason *regexp.Regexp
	User   *regexp.Regexp
}

var AuthRegexInstance = authRegex{
	Ip:     regexp.MustCompile(`[from ]\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`),
	Port:   regexp.MustCompile(`port \d{2,}`),
	Reason: regexp.MustCompile(`\d:.*`),
	User:   regexp.MustCompile(`user [[:alpha:]].* `),
}

type authParser struct {
	AuthEvent
}

var _ AuthParser = (*authParser)(nil)

var AuthParserInstance = authParser{}

// [parseAuthTS] parses nginx timestamp string
// using layout "2006-01-02T15:04:05.000000+00:00"
// and converts it to a string of the layout "2006-01-02 15:04:05"
func parseAuthTS(s string) string {
	const layout = "2006-01-02T15:04:05.000000+00:00"
	const target_layout = "2006-01-02 15:04:05"
	parsed, _ := time.Parse(layout, s)
	formatted := parsed.Format(target_layout)
	return formatted
}

func (a *authParser) AuthParse(line string, out chan<- AuthEvent) {
	if !(strings.Contains(line, "ndudaev") ||
		strings.Contains(line, "pam_unix") ||
		strings.Contains(line, "lastlog_openseek") ||
		strings.Contains(line, "new user")) {
		authevent := AuthEvent{}
		ts := strings.SplitN(line, " ", 2)
		authevent.Ts = parseAuthTS(ts[0])
		skipMeta := strings.SplitN(line, "]:", 2)
		connectionDetails := skipMeta[1]
		re := AuthRegexInstance.Ip.Find([]byte(connectionDetails))
		authevent.Ip = strings.TrimSpace(string(re))
		re = AuthRegexInstance.Port.Find([]byte(connectionDetails))
		replacer := strings.NewReplacer("port ", "")
		authevent.Port = replacer.Replace(strings.TrimSpace(string(re)))
		re = AuthRegexInstance.Reason.Find([]byte(connectionDetails))
		authevent.Reason = string(re)
		re = AuthRegexInstance.User.Find([]byte(connectionDetails))
		found_str := strings.SplitN(string(re), " ", 3)
		if len(found_str) > 2 {
			authevent.User = found_str[1]
		}
		if strings.Contains(connectionDetails, "user") {
			action, _, _ := strings.Cut(connectionDetails, "user ")
			authevent.Action = strings.TrimSpace(action)
		} else {
			action, _, _ := strings.Cut(connectionDetails, authevent.Ip)
			authevent.Action = strings.TrimSpace(action)
		}

		out <- authevent
	}
}
