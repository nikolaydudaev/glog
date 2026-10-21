package parse

import (
	"encoding/hex"
	"log"
	"strconv"
	"strings"
	"time"
)

type ufwParser struct {
	FirewallEvent
}

var _ FirewallParser = (*ufwParser)(nil)

var UfwParserInstance = ufwParser{}

// [parseUfwTS] parses nginx timestamp string
// using layout "2006-01-02T15:04:05.000000+00:00"
// and converts it to a string of the layout "2006-01-02 15:04:05"
func parseUfwTS(s string) string {
	const layout = "2006-01-02T15:04:05.000000+00:00"
	const target_layout = "2006-01-02 15:04:05"
	parsed, _ := time.Parse(layout, s)
	formatted := parsed.Format(target_layout)
	return formatted
}

func (f *ufwParser) FirewallParse(line string, out chan<- FirewallEvent) {
	fevent := FirewallEvent{}
	splitLine := strings.SplitN(line, " ", 3)
	fevent.Ts = parseUfwTS(splitLine[0])
	fevent.VpsName = splitLine[1]
	splitLine = strings.SplitN(line, "[", 2)
	splitLine = strings.SplitN(splitLine[1], "]", 2)
	fevent.UFWSignal = splitLine[0]
	splitLine = strings.Split(splitLine[1], " ")
	for _, item := range splitLine {
		if strings.Contains(item, "=") {
			itemSplit := strings.Split(item, "=")
			switch itemSplit[0] {
			case "IN":
				if len(itemSplit) == 1 {
					fevent.InputInterface = ""
				} else {
					fevent.InputInterface = itemSplit[1]
				}
			case "OUT":
				if len(itemSplit) == 1 {
					fevent.OutputInterface = ""
				} else {
					fevent.OutputInterface = itemSplit[1]
				}
			case "MAC":
				if len(itemSplit) == 1 {
					fevent.MACAddress = ""
				} else {
					fevent.MACAddress = itemSplit[1]
				}
			case "SRC":
				if len(itemSplit) == 1 {
					fevent.SourceIp = ""
				} else {
					fevent.SourceIp = itemSplit[1]
				}
			case "DST":
				if len(itemSplit) == 1 {
					fevent.DestIp = ""
				} else {
					fevent.DestIp = itemSplit[1]
				}
			case "LEN":
				if len(itemSplit) == 1 {
					fevent.Length = nil
				} else {
					n, err := strconv.Atoi(itemSplit[1])
					if err != nil {
						log.Fatal(err)
					} else {
						fevent.Length = &n
					}
				}
			case "TOS":
				if len(itemSplit) == 1 {
					fevent.TypeOfService = []byte{}
				} else {
					fevent.TypeOfService, _ = hex.DecodeString(itemSplit[1])
				}
			case "PREC":
				if len(itemSplit) == 1 {
					fevent.PrecedenceBit = []byte{}
				} else {
					fevent.PrecedenceBit, _ = hex.DecodeString(itemSplit[1])
				}
			case "TTL":
				if len(itemSplit) == 1 {
					fevent.TimeToLive = nil
				} else {
					n, err := strconv.Atoi(itemSplit[1])
					if err != nil {
						log.Fatal(err)
					} else {
						fevent.TimeToLive = &n
					}
				}
			case "ID":
				if len(itemSplit) == 1 {
					fevent.Id = nil
				} else {
					n, err := strconv.Atoi(itemSplit[1])
					if err != nil {
						log.Fatal(err)
					} else {
						fevent.Id = &n
					}
				}
			case "PROTO":
				if len(itemSplit) == 1 {
					fevent.Proto = ""
				} else {
					fevent.Proto = itemSplit[1]
				}
			case "SPT":
				if len(itemSplit) == 1 {
					fevent.SourcePort = ""
				} else {
					fevent.SourcePort = itemSplit[1]
				}
			case "DPT":
				if len(itemSplit) == 1 {
					fevent.DestPort = ""
				} else {
					fevent.DestPort = itemSplit[1]
				}
			case "WINDOW":
				if len(itemSplit) == 1 {
					fevent.Window = nil
				} else {
					n, err := strconv.Atoi(itemSplit[1])
					if err != nil {
						log.Fatal(err)
					} else {
						fevent.Window = &n
					}
				}
			case "RES":
				if len(itemSplit) == 1 {
					fevent.ReservedBits = []byte{}
				} else {
					fevent.ReservedBits, _ = hex.DecodeString(itemSplit[1])
				}
			case "URGP":
				if len(itemSplit) == 1 {
					fevent.UrgentPointer = nil
				} else {
					n, err := strconv.Atoi(itemSplit[1])
					if err != nil {
						log.Fatal(err)
					} else {
						fevent.UrgentPointer = &n
					}
				}
			case "TYPE":
				if len(itemSplit) == 1 {
					fevent.ICMPType = nil
				} else {
					n, err := strconv.Atoi(itemSplit[1])
					if err != nil {
						log.Fatal(err)
					} else {
						fevent.ICMPType = &n
					}
				}
			case "CODE":
				if len(itemSplit) == 1 {
					fevent.ICMPCode = nil
				} else {
					n, err := strconv.Atoi(itemSplit[1])
					if err != nil {
						log.Fatal(err)
					} else {
						fevent.ICMPCode = &n
					}
				}
			case "MARK":
				if len(itemSplit) == 1 {
					fevent.PacketMark = ""
				} else {
					fevent.PacketMark = itemSplit[1]
				}
			}
		} else {
			switch item {
			case "DF", "MF", "CE":
				fevent.ECNIndicators = append(fevent.ECNIndicators, item)
			case "SYN", "ACK", "FIN", "RST", "PSH", "URG", "ECE", "CWR":
				fevent.TCPFlags = append(fevent.TCPFlags, item)
			}
		}
	}

	out <- fevent
}
