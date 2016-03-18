package ais

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	//"fmt"
	"strconv"
	"strings"
)

type parser_item struct {
	MMSI             uint32
	Message_type     uint8
	Repeat_indicator uint8
	Binary           []byte
	Valid            bool
}

type Parser struct {
	fragments map[uint8]parser_item
}

func GetParser() Parser {
	parser := Parser{}
	parser.fragments = make(map[uint8]parser_item) // [message_id][]raw_string
	return parser
}

/*
	парсим исходную строку
*/
func (this *Parser) Parse(data string) (AIS, bool) {
	var (
		split []string = strings.Split(data, ",")
		body  string   = split[5]
		runes []rune
		err   error
		value uint64
	)

	var (
		count_fragments uint8 // количество фрагментов
		fragment_number uint8 // номер фрагмента
		message_id      uint8 // идентификатор фрагментов
	)

	if len(split) < 7 {
		return AIS{}, true
	}

	//count_fragments
	value, err = strconv.ParseUint(split[1], 10, 8)
	if err != nil {
		return AIS{}, true
	}
	count_fragments = uint8(value)

	//fragment_number
	value, err = strconv.ParseUint(split[2], 10, 8)
	if err != nil {
		return AIS{}, true
	}
	fragment_number = uint8(value)

	//message_id
	if split[3] != "" {
		value, err = strconv.ParseUint(split[3], 10, 8)
		if err != nil {
			return AIS{}, true
		}
		message_id = uint8(value)
	}

	for i := 0; i < len(body); i++ {
		runes = append(runes, ascii_to_8bit(rune(body[i])))
	}

	bits := dec_to_6bit(runes)

	var res_bits []byte
	var item_ais AIS

	if count_fragments == 1 || fragment_number == 1 {
		// Message Type
		res_bits = getBinary(bits, 0, 6)
		message_type := uint8(res_bits[0])

		// Repeat Indicator
		res_bits = getBinary(bits, 6, 2)
		repeat_indicator := uint8(res_bits[0])

		// MMSI
		res_bits = getBinary(bits, 8, 30)
		mmsi := binary.BigEndian.Uint32(res_bits)

		if count_fragments > 1 {
			this.fragments[message_id] = parser_item{
				MMSI:             mmsi,
				Message_type:     message_type,
				Repeat_indicator: repeat_indicator,
				Binary:           bits,
				Valid:            this.checkCRC(data) == nil,
			}

			return AIS{}, false // сообщаем, что объект ещё не создан
		} else {
			// одиночное сообщение
			item_ais = this.createAIS(mmsi, split[4], message_type, repeat_indicator, this.checkCRC(data) == nil, bits)
			return item_ais, true
		}
	} else {
		// остальные фрагменты
		var (
			fragment parser_item
			ok       bool
		)
		if fragment, ok = this.fragments[message_id]; !ok {
			return AIS{}, true
		}
		//fmt.Println(bits)

		fragment.Binary = append(fragment.Binary, bits...)
		fragment.Valid = fragment.Valid && this.checkCRC(data) == nil

		//this.fragments[message_id].Binary = append(this.fragments[message_id].Binary, bits...)

		// если последнее сообщение
		if count_fragments == fragment_number {
			item_ais = this.createAIS(fragment.MMSI, split[4], fragment.Message_type, fragment.Repeat_indicator, fragment.Valid, fragment.Binary)

			delete(this.fragments, message_id)

			return item_ais, true
		} else {
			this.fragments[message_id] = fragment
			return AIS{}, false // сообщаем, что объект ещё не создан
		}
	}

	return AIS{}, false
}

/*
	создание объекта AIS
*/
func (this Parser) createAIS(mmsi uint32, radio_channel string, message_type, repeat_indicator uint8, valid bool, bits []byte) AIS {
	var item_ais AIS = AIS{
		MMSI:          mmsi,
		Message_type:  message_type,
		Radio_channel: radio_channel,
		Valid:         valid,
		is_parsed:     true,
	}

	switch message_type {
	case 1, 2, 3:
		this.parsePositionReport(&item_ais, bits)
	case 4:
		this.parseBaseStationReport(&item_ais, bits)
	case 5:
		this.parseShipAndVoyage(&item_ais, bits)
	}

	return item_ais
}

func (this Parser) checkCRC(data string) error {
	var (
		body    string
		crc     int = 0
		index_s int = strings.Index(data, "*")
	)

	if strings.Contains(data, "!") && strings.Contains(data, "*") {
		body = data[1:index_s]
	} else {
		body = data
	}

	for i := 0; i < len(body); i++ {
		if crc == 0 {
			crc = int(rune(body[i]))
		} else {
			crc ^= int(rune(body[i]))
		}
	}

	if strings.ToUpper(hex.EncodeToString([]byte{byte(crc)})) != strings.ToUpper(data[index_s+1:index_s+3]) {
		return errors.New("CRC error")
	}

	return nil
}
