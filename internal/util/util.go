package util

import "strings"

// GetHexString converts a byte array like:
// {0x04, 0xA7, 0x0F, 0x13} to a string like:
// "04 A7 0F 13"
func GetHexString(hexBytes []byte) string {
	const hextable = "0123456789ABCDEF"
	dst := make([]byte, len(hexBytes)*3)
	j := 0
	for _, v := range hexBytes {
		dst[j] = hextable[v>>4]
		dst[j+1] = hextable[v&0x0f]
		dst[j+2] = byte(' ')
		j += 3
	}
	return strings.TrimSpace(string(dst))
}

// getPosHex convers a byte array that looks like:
// {0x04, 0x07, 0x0F, 0x03} to a string like: 0x47F3
func GetPosHex(value []byte) string {
	const hextable = "0123456789ABCDEF"
	hexString := "0x"
	for _, b := range value {
		char := string(hextable[b&0x0F])
		hexString = hexString + char
	}
	return hexString
}

// getBytesFromPosHex converts a 4 digit hex number like 0x47F3
// to a byte array like this: {0x04, 0x07, 0x0F, 0x03}
func GetBytesFromPosHex(value int) []byte {
	byte1 := byte(value & 0xF000 >> 12)
	byte2 := byte(value & 0x0F00 >> 8)
	byte3 := byte(value & 0x00F0 >> 4)
	byte4 := byte(value & 0x000F)
	return []byte{byte1, byte2, byte3, byte4}
}
