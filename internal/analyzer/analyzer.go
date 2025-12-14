package analyzer

import (
	"strings"

	// "github.com/google/gopacket"
	"github.com/goushalk/lazyshark/internal/pcapreader"
)

type PacketSummary struct {
	Number    int
	TimeStamp string
	SrcIp     string
	DstIp     string
	Protocol  string
	Length    int
	Info      string   // SAFE summary only
	RawData   []byte   // FULL packet bytes
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

		protocol := "UNKNOWN"
		if pkt.TransportLayer() != nil {
			protocol = pkt.TransportLayer().LayerType().String()
		}

		info := "no application data"

		if app := pkt.ApplicationLayer(); app != nil {
			data := app.Payload()

			if len(data) > 0 {
				// VERY conservative text detection
				if strings.HasPrefix(string(data), "GET") ||
					strings.HasPrefix(string(data), "POST") ||
					strings.HasPrefix(string(data), "HTTP") {
					info = "HTTP"
					protocol = "HTTP"
				} else {
					info = "binary payload"
				}
			}
		}

		src, dst := "N/A", "N/A"
		if net := pkt.NetworkLayer(); net != nil {
			src = net.NetworkFlow().Src().String()
			dst = net.NetworkFlow().Dst().String()
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
