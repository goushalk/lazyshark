package pcapreader

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

func Reader(filePath string) ([]gopacket.Packet, error) {
	// Expand ~ manually

	file := filePath
	handle, err := pcap.OpenOffline(file)
	if err != nil {
		return nil, err
	}
	defer handle.Close()

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	var packets []gopacket.Packet
	for packet := range packetSource.Packets() {
		packets = append(packets, packet)
	}
	return packets, nil
}
