package maxim

const crcInit uint32 = 0x07 // 5-bit init word, constant, 00111
const crcPoly uint32 = 0x35 // 6-bit polynomial, constant, 110101
const crcLength int = 19    // 19-bit data

// / <summary>
// / Last 5 bits of 24-bit frame are discarded, 19-bit data length
// / Polynomial P(x) = x5+x4+x2+x0 -> 110101
// / Init word append to the 19-bit data -> 00111
// / </summary>
// / <param name="data2">MSB Byte </param>
// / <param name="data1">Middle Byte </param>
// / <param name="data0">LSB Byte </param>
// / <returns>5-bit CRC </returns>
func CRC(data []byte) (byte, error) {

	var crcStep uint32 = 0
	var crcResult byte = 0
	var tmp uint32 = 0

	datainput := uint32(data[0])<<16 + uint32(data[1])<<8 + uint32(data[2]&0xe0)

	//append 5-bit init word to first 19-bit data
	datainput = (datainput & 0xffffe0) + crcInit

	//first step, get crc_step 0
	tmp = ((datainput & 0xfc0000) >> 18) //crc_step 0= data[18:13]
	//next crc_step = crc_step[5] = 0 ? (crc_step[5:0] ^ crc_poly) : crc_step[5:0]
	if (tmp & 0x20) == 0x20 {
		crcStep = (tmp ^ crcPoly)
	} else {
		crcStep = tmp
	}
	//step 1-18
	for i := 0; i < crcLength-1; i++ {
		//append next data bit to previous crc_step[4:0], {crc_step[4:0], next data bit}
		tmp = (((crcStep & 0x1f) << 1) + ((datainput >> (crcLength - 2 - i)) & 0x01))
		//next crc_step = crc_step[5] = 0 ? (crc_step[5:0] ^ crc_poly) : crc_step[5:0]
		if (tmp & 0x20) == 0x20 {
			crcStep = (tmp ^ crcPoly)
		} else {
			crcStep = tmp
		}
	}

	crcResult = byte(crcStep & 0x1f) //crc result = crc_step[4:0]
	return crcResult, nil
}
