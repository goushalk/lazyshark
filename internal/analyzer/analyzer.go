package analyzer

import (
	"strings"

	"github.com/google/gopacket/layers"

	"packetB/internal/pcapreader"
)

type PacketSummary struct {
	Number    int
	TimeStamp string
	SrcIp     string
	DstIp     string
	Protocol  string
	Length    int
	Info      string
	RawData   []byte
}

type AnalyzerResult struct {
	Packets        []PacketSummary
	ProtocolCounts map[string]int
}

func Analyzer(file string) (*AnalyzerResult, error) {
	packets, err := pcapreader.Reader(file)
	if err != nil {
		return nil, err
	}

	result := &AnalyzerResult{
		Packets:        make([]PacketSummary, 0, len(packets)),
		ProtocolCounts: make(map[string]int),
	}

	for i, pkt := range packets {

		// ---------------------------
		// Network layer (IPv4 / IPv6)
		// ---------------------------
		src, dst := "N/A", "N/A"

		if ip4 := pkt.Layer(layers.LayerTypeIPv4); ip4 != nil {
			ip := ip4.(*layers.IPv4)
			src = ip.SrcIP.String()
			dst = ip.DstIP.String()
		} else if ip6 := pkt.Layer(layers.LayerTypeIPv6); ip6 != nil {
			ip := ip6.(*layers.IPv6)
			src = ip.SrcIP.String()
			dst = ip.DstIP.String()
		}

		// ---------------------------
		// Transport layer
		// ---------------------------
		transport := "UNKNOWN"
		if tl := pkt.TransportLayer(); tl != nil {
			transport = tl.LayerType().String()
		}

		// ---------------------------
		// Application / Info
		// ---------------------------
		info := "no application data"
		appProto := ""

		if app := pkt.ApplicationLayer(); app != nil {
			data := app.Payload()

			if len(data) > 0 {
				// Safe HTTP detection only
				if strings.HasPrefix(string(data), "GET ") ||
					strings.HasPrefix(string(data), "POST ") ||
					strings.HasPrefix(string(data), "HTTP/") {
					appProto = "HTTP"
					info = "HTTP request/response"
				} else {
					info = "binary payload"
				}
			}
		}

		protocol := transport
		if appProto != "" {
			protocol = transport + "/" + appProto
		}

		summary := PacketSummary{
			Number:    i + 1,
			TimeStamp: pkt.Metadata().Timestamp.String(),
			SrcIp:     src,
			DstIp:     dst,
			Protocol:  protocol,
			Length:    pkt.Metadata().Length,
			Info:      info,
			RawData:   pkt.Data(),
		}

		result.Packets = append(result.Packets, summary)
		result.ProtocolCounts[protocol]++
	}

	return result, nil
}
