package ais

import "math"

type Bits_array []bool

func (this Bits_array) GetBinary(start_bit, length int) []byte {
	var (
		length_responce int    = int(math.Ceil(float64(length) / 8))
		responce        []byte = make([]byte, length_responce)
		slice           []bool = this[start_bit : start_bit+length]
		pos_byte_res    int    = length_responce - 1
		pos_bit_res     byte   = 0x00
	)

	for i := len(slice) - 1; i > -1; i-- {
		if slice[i] {
			responce[pos_byte_res] |= 0x01 << pos_bit_res
		}

		pos_bit_res++
		if pos_bit_res == 0x08 {
			pos_byte_res--
			pos_bit_res = 0x00
		}
	}

	return responce
}

func (this Bits_array) GetText(start_bit, length int) string {
	var (
		responce    []byte
		slice       []bool = this[start_bit : start_bit+length]
		pos_bit_res byte   = 0x05
		item_byte   byte   = 0x00
	)

	for i := 0; i < len(slice); i++ {
		if slice[i] {
			item_byte |= 0x01 << pos_bit_res
		}

		if pos_bit_res == 0x00 {
			if item_byte > 0x00 {
				responce = append(responce, byte8bit_to_ascii(item_byte))
			}

			item_byte = 0x00
			pos_bit_res = 0x05
			continue
		}

		pos_bit_res--
	}

	return string(responce)
}
