package ais

import (
	"bytes"
	"math"
)

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

func dec_to_6bit(data []rune) []byte {
	var (
		b         byte
		item_byte byte
		pos_bit   byte = 0x00
		buffer    bytes.Buffer
	)

	for i := 0; i < len(data); i++ {
		for b = 6; b > 0; b-- {
			if (data[i]>>(b-1))&0x01 == 0x01 {
				item_byte |= 0x01 << (0x07 - pos_bit)
			}

			pos_bit++
			if pos_bit == 0x08 {
				pos_bit = 0x00
				buffer.WriteByte(item_byte)
				item_byte = 0x00
			}
		}
	}

	if pos_bit != 0x08 {
		buffer.WriteByte(item_byte)
	}

	return buffer.Bytes()
}

/*
	получить байты в диапазоне бит
*/
func getBinary(body []byte, start_bit, length int) []byte {
	var (
		start_byte      int    = int(math.Ceil(float64(start_bit+1)/8)) - 1    // порядковый номер первого байта
		length_body     int    = int(math.Ceil(float64(start_bit+length) / 8)) // количество байт нужно проверить
		length_responce int    = int(math.Ceil(float64(length) / 8))           // длина ответа
		responce        []byte = make([]byte, length_responce)
	)

	var (
		pos_byte_res int  = length_responce - 1
		pos_bit_res  byte = 0x00                                               // позиция бита в байте ответа
		pos_bit_body byte = ((byte(length_body) * 8) - byte(start_bit+length)) // позиция бита в байте входящих данных
		i            byte
	)

	for pos := length_body - 1; pos >= start_byte; pos-- {
		for i = pos_bit_body; i < 0x08; i++ {
			// если бит включен
			if (body[pos]>>i)&0x01 == 0x01 {
				responce[pos_byte_res] |= 0x01 << pos_bit_res
			}

			// выходим из цикла, когда обошли все биты
			length--
			if length == 0 {
				goto Exit
			}

			pos_bit_res++
			if pos_bit_res == 0x08 {
				// если последний байт уже заполнен, то выходим
				if pos_byte_res == 0x00 {
					goto Exit
				}
				pos_byte_res--
				pos_bit_res = 0x00
			}
		}

		pos_bit_body = 0x00
	}

Exit:

	return responce
}
