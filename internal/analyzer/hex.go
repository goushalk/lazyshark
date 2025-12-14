
package analyzer

import "encoding/hex"

func DumpHex(data []byte) string {
	if len(data) == 0 {
		return "no packet data"
	}
	return hex.Dump(data)
}
