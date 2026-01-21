package pcapreader

import (
	"io"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func Reader(filePath string) ([]gopacket.Packet, error) {
	handle, err := pcap.OpenOffline(filePath)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	var packets []gopacket.Packet
	linkType := handle.LinkType()

	for {
		data, ci, err := handle.ReadPacketData()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		var firstLayer gopacket.LayerType

		switch linkType {

		case layers.LinkTypeEthernet:
			firstLayer = layers.LayerTypeEthernet

		case layers.LinkTypeLinuxSLL:
			firstLayer = layers.LayerTypeLinuxSLL

		case layers.LinkTypeLoop:
			firstLayer = layers.LayerTypeLoopback

		case layers.LinkTypeRaw, layers.LinkTypeIPv4:
			// 🔥 THIS IS YOUR FILE
			firstLayer = layers.LayerTypeIPv4

		default:
			// Last-resort fallback
			firstLayer = layers.LayerTypeIPv4
		}

		packet := gopacket.NewPacket(
			data,
			firstLayer,
			gopacket.DecodeOptions{
				Lazy:   false,
				NoCopy: false,
			},
		)

		packet.Metadata().CaptureInfo = ci
		packets = append(packets, packet)
	}

	return packets, nil
}
