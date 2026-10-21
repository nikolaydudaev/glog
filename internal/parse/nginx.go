package parse

import (
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type nginxParser struct {
	source        string
	client_ip     *regexp.Regexp
	ts            *regexp.Regexp
	request       *regexp.Regexp
	http_status   *regexp.Regexp
	response_size *regexp.Regexp
	host          *regexp.Regexp
	user_agent    *regexp.Regexp
}

var _ Parser = (*nginxParser)(nil)

var NginxParserInstance = nginxParser{
	source:        "nginx",
	client_ip:     regexp.MustCompile(`^\b((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`),
	ts:            regexp.MustCompile(`\[\d{2}\/\w{3}\/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4}\]`),
	request:       regexp.MustCompile(`"[A-Z].+ \/.* (HTTP)\/\d{1}\.\d{1}"`),
	http_status:   regexp.MustCompile(` \d{1,} `),
	response_size: regexp.MustCompile(` \d{1,}`),
	host:          regexp.MustCompile(` "(-|[https].*)" `),
}

func findToken(s []byte, re *regexp.Regexp, field string) string {
	var found_str string
	switch field {
	case "http_status":
		if len(re.FindAll([]byte(s), 1)) > 0 {
			found_str = string(re.FindAll([]byte(s), 1)[0])
		} else {
			found_str = ""
		}
	case "response_size":
		if len(re.FindAll([]byte(s), 2)) > 2 {
			found_str = string(re.FindAll([]byte(s), 2)[1])
		} else {
			found_str = ""
		}

	case "user_agent":
		found_str = strings.Split(string(s), ` "`)[len(strings.Split(string(s), ` "`))-1]
	default:
		found_str = string(re.Find([]byte(s)))
	}
	found_str = strings.TrimSpace(found_str)
	if field == "ts" {
		replacer := strings.NewReplacer("[", "", "]", "")
		found_str = replacer.Replace(found_str)
	}
	if field == "user_agent" || field == "host" || field == "request" {
		replacer := strings.NewReplacer(`"`, "", "-", "")
		found_str = replacer.Replace(found_str)
	}

	return found_str
}

// [parseTS] parses nginx timestamp string
// using layout "02/Jan/2006:15:04:05 -0700"
// and converts it to a string of the layout "2006-01-02 15:04:05"
func parseTS(s string) string {
	const layout = "02/Jan/2006:15:04:05 -0700"
	const target_layout = "2006-01-02 15:04:05"
	parsed, _ := time.Parse(layout, s)
	formatted := parsed.Format(target_layout)
	return formatted
}

// [extractMethod] accepts a request string
// It cuts the string by space and pulls the method name from the string
func extractMethod(req string) string {
	method, _, _ := strings.Cut(req, " ")
	return method
}

// [extractProto] accepts a request string
// It extracts the protocol name and version
func extractProto(req string) string {
	re := regexp.MustCompile(`(HTTP/).*$`)
	found_str := string(re.Find([]byte(req)))
	found_str = strings.TrimSpace(found_str)
	return found_str
}

// [extractURI] attepmts to extract requests uri
func extractURI(req string) string {
	re := regexp.MustCompile(` /.* | \* `)
	found_str := string(re.Find([]byte(req)))
	found_str = strings.TrimSpace(found_str)
	return found_str
}

func (p *nginxParser) Parse(line string, out chan<- ServerEvent) {
	event := ServerEvent{}
	event.Source = p.source
	event.Msg = ""
	event.RemoteIp = ""
	event.RemotePort = ""
	event.RequestHeaders = http.Header{}
	event.ServerName = ""
	event.BytesRead = 0
	event.Duration = 0.0
	event.ClientIp = findToken([]byte(line), p.client_ip, "client_ip")
	event.Ts = parseTS(findToken([]byte(line), p.ts, "ts"))
	pre_request := findToken([]byte(line), p.request, "request")
	event.Method = extractMethod(pre_request)
	event.Proto = extractProto(pre_request)
	event.Uri = extractURI(pre_request)
	event.Status, _ = strconv.Atoi(findToken([]byte(line), p.http_status, "http_status"))
	event.Size, _ = strconv.Atoi(findToken([]byte(line), p.response_size, "response_size"))
	event.Host = findToken([]byte(line), p.host, "host")
	event.RequestHeaders.Set("User-Agent", findToken([]byte(line), p.user_agent, "user_agent"))

	out <- event
}
