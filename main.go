package ais

func ascii_to_8bit(ascii rune) rune {
	if ascii > 119 {
		return ascii
	}

	if ascii > 87 && ascii < 96 {
		return ascii
	}

	ascii += 40
	if ascii > 128 {
		ascii += 32
	} else {
		ascii += 40
	}

	return ascii
}

func byte8bit_to_ascii(char byte) byte {
	if char < 32 {
		char += 64
	}

	return char
}

func dec_to_6bit(data []rune) Bits_array {
	var (
		b    byte
		bits Bits_array
	)

	for i := 0; i < len(data); i++ {
		for b = 0x06; b > 0x00; b-- {
			//for b = 0x01; b < 0x07; b++ {
			bits = append(bits, (data[i]>>(b-1))&0x01 == 0x01)
		}
	}

	return bits
}
