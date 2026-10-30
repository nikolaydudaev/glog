package parse

import (
	"net/http"
	"time"
)

type Parser interface {
	Parse(line string, out chan<- ServerEvent)
}

type FirewallParser interface {
	FirewallParse(line string, out chan<- FirewallEvent)
}

type AuthParser interface {
	AuthParse(line string, out chan<- AuthEvent)
}

type FailbanParser interface {
	FailbanParse(line string, out chan<- FailbanEvent)
}

type ServerEvent struct {
	Source         string
	Level          string
	Logger         string
	Ts             string //time.Time // nginx needs date parsing and converting; caddy needs converting from UTC
	Msg            string // only caddy, empty for nginx
	RemoteIp       string // only caddy, empty for nginx
	RemotePort     string // only caddy, empty for nginx
	ClientIp       string
	Proto          string      // nginx needs extracting from request
	Method         string      // nginx needs extracting from request
	Host           string      // it is page in old nginx logs
	Uri            string      // nginx needs extracting from request
	RequestHeaders http.Header //only caddy
	TlsResumed     bool
	TlsVersion     uint16
	TlsCipherSuite uint16
	TlsProto       string
	TlsServerName  string // from tls section of caddy logs
	TlsEch         bool
	BytesRead      int           // only caddy
	Duration       time.Duration //only caddy
	Size           int           // response size in nginx
	Status         int
}

type FirewallEvent struct {
	Ts              string
	VpsName         string
	UFWSignal       string
	InputInterface  string
	OutputInterface string
	MACAddress      string
	SourceIp        string
	SourcePort      string
	DestIp          string
	DestPort        string
	Length          *int
	TypeOfService   []byte
	PrecedenceBit   []byte
	TimeToLive      *int
	Id              *int
	ECNIndicators   []string
	Proto           string
	Window          *int
	ReservedBits    []byte
	TCPFlags        []string
	UrgentPointer   *int
	ICMPType        *int
	ICMPCode        *int
	PacketMark      string
}

type AuthEvent struct {
	Ts     string
	Ip     string
	Port   string
	Reason string
	User   string
	Action string
}

type FailbanEvent struct {
	Ts      string
	Flag    string
	Message string
	Action  string
	Ip      string
}
