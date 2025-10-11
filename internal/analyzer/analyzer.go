package analyzer

import (
	"strings"

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

		protocol := "UNKNOWN"
		if pkt.TransportLayer() != nil {
			protocol = pkt.TransportLayer().LayerType().String()
		}

		payload := "no data"
		if pkt.ApplicationLayer() != nil {
			tem_payload := string(pkt.ApplicationLayer().Payload())
			if strings.HasPrefix(tem_payload, "HTTP") || strings.Contains(tem_payload, "GET") || strings.Contains(tem_payload, "POST") {
				payload = "HTTP"
				protocol = "HTTP"
			} else {
				payload = tem_payload
			}
		}

		networkLayer := pkt.NetworkLayer()
		if networkLayer == nil {
			summary := PacketSummary{
				Number:    i + 1,
				TimeStamp: pkt.Metadata().Timestamp.String(),
				SrcIp:     "N/A",
				DstIp:     "N/A",
				Protocol:  protocol,
				Length:    pkt.Metadata().Length,
				Info:      payload,
				RawData:   pkt.Data(),
			}
			result.Packets = append(result.Packets, summary)
			result.ProtocolCounts[protocol]++
			continue

		}

		summary := PacketSummary{
			Number:    i + 1,
			TimeStamp: pkt.Metadata().Timestamp.String(),
			SrcIp:     networkLayer.NetworkFlow().Src().String(),
			DstIp:     networkLayer.NetworkFlow().Dst().String(),
			Protocol:  protocol,
			Length:    pkt.Metadata().Length,
			Info:      payload,
			RawData:   pkt.Data(),
		}

		result.Packets = append(result.Packets, summary)
		result.ProtocolCounts[summary.Protocol]++
	}
	return result, nil

}
