package parse

import (
	"encoding/json/v2"
	"log"
	"math"
	"net/http"
	"time"
)

type TLS struct {
	TlsResumed     bool   `json:"resumed"`
	TlsVersion     uint16 `json:"version"`
	TlsCipherSuite uint16 `json:"cipher_suite"`
	TlsProto       string `json:"proto"`
	TlsServerName  string `json:"server_name"`
	TlsEch         bool   `json:"ech"`
}

type Request struct {
	RemoteIp   string      `json:"remote_ip"`
	RemotePort string      `json:"remote_port"`
	ClientIp   string      `json:"client_ip"`
	Proto      string      `json:"proto"`
	Method     string      `json:"method"`
	Host       string      `json:"host"`
	Uri        string      `json:"uri"`
	Headers    http.Header `json:"headers"`
	TLS        TLS         `json:"tls"`
}

type CaddyLine struct {
	Level       string      `json:"level"`
	Ts          float64     `json:"ts"`
	Logger      string      `json:"logger"`
	Msg         string      `json:"msg"`
	Request     Request     `json:"request"`
	BytesRead   int         `json:"bytes_read"`
	UserId      string      `json:"user_id"`
	Duration    float64     `json:"duration"`
	Size        int         `json:"size"`
	Status      int         `json:"status"`
	RespHeaders http.Header `json:"resp_headers"`
}

type caddyParser struct {
	CaddyLine
}

var _ Parser = (*caddyParser)(nil)

var CaddyParserInstance = caddyParser{}

func parseUnixTS(f float64) string {
	sec, frac := math.Modf(f)
	t := time.Unix(int64(sec), int64(frac*1e9)).UTC()
	const target_layout = "2006-01-02 15:04:05"
	formatted := t.Format(target_layout)
	return formatted
}

func (p *caddyParser) Parse(line string, out chan<- ServerEvent) {
	var caddyLine CaddyLine
	var event ServerEvent
	err := json.Unmarshal([]byte(line), &caddyLine)
	if err != nil {
		log.Fatal(err)
	}

	event.Source = "caddy"
	event.Ts = parseUnixTS(caddyLine.Ts)
	event.Level = caddyLine.Level
	event.Logger = caddyLine.Logger
	event.Msg = caddyLine.Msg
	event.RemoteIp = caddyLine.Request.RemoteIp
	event.RemotePort = caddyLine.Request.RemotePort
	event.ClientIp = caddyLine.Request.ClientIp
	event.Proto = caddyLine.Request.Proto
	event.Method = caddyLine.Request.Method
	event.Host = caddyLine.Request.Host
	event.Uri = caddyLine.Request.Uri
	event.RequestHeaders = caddyLine.Request.Headers
	event.TlsResumed = caddyLine.Request.TLS.TlsResumed
	event.TlsVersion = caddyLine.Request.TLS.TlsVersion
	event.TlsCipherSuite = caddyLine.Request.TLS.TlsCipherSuite
	event.TlsProto = caddyLine.Request.TLS.TlsProto
	event.TlsServerName = caddyLine.Request.TLS.TlsServerName
	event.TlsEch = caddyLine.Request.TLS.TlsEch
	event.BytesRead = caddyLine.BytesRead
	event.Duration = time.Duration(caddyLine.Duration * float64(time.Second))
	event.Size = caddyLine.Size
	event.Status = caddyLine.Status

	out <- event
}
