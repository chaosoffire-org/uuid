package internal

import "crypto/rand"

// encodeHex encodes uuid into dst using branchless bit manipulation.
func encodeHex(dst *[36]byte, uuid UUID) {
	dst[8] = '-'
	dst[13] = '-'
	dst[18] = '-'
	dst[23] = '-'

	dst[35] = hextable[uuid[15]&0x0f]
	dst[34] = hextable[uuid[15]>>4]
	dst[33] = hextable[uuid[14]&0x0f]
	dst[32] = hextable[uuid[14]>>4]
	dst[31] = hextable[uuid[13]&0x0f]
	dst[30] = hextable[uuid[13]>>4]
	dst[29] = hextable[uuid[12]&0x0f]
	dst[28] = hextable[uuid[12]>>4]
	dst[27] = hextable[uuid[11]&0x0f]
	dst[26] = hextable[uuid[11]>>4]
	dst[25] = hextable[uuid[10]&0x0f]
	dst[24] = hextable[uuid[10]>>4]

	dst[22] = hextable[uuid[9]&0x0f]
	dst[21] = hextable[uuid[9]>>4]
	dst[20] = hextable[uuid[8]&0x0f]
	dst[19] = hextable[uuid[8]>>4]

	dst[17] = hextable[uuid[7]&0x0f]
	dst[16] = hextable[uuid[7]>>4]
	dst[15] = hextable[uuid[6]&0x0f]
	dst[14] = hextable[uuid[6]>>4]

	dst[12] = hextable[uuid[5]&0x0f]
	dst[11] = hextable[uuid[5]>>4]
	dst[10] = hextable[uuid[4]&0x0f]
	dst[9] = hextable[uuid[4]>>4]

	dst[7] = hextable[uuid[3]&0x0f]
	dst[6] = hextable[uuid[3]>>4]
	dst[5] = hextable[uuid[2]&0x0f]
	dst[4] = hextable[uuid[2]>>4]
	dst[3] = hextable[uuid[1]&0x0f]
	dst[2] = hextable[uuid[1]>>4]
	dst[1] = hextable[uuid[0]&0x0f]
	dst[0] = hextable[uuid[0]>>4]
}

func readRandom(dst *[16]byte) error {
	_, err := rand.Read(dst[:])
	if err != nil {
		return err
	}

	return nil
}
