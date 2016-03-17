package ais

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"strconv"
	"strings"
)

type AIS struct {
	Static
	Dynamic
	last_measure
}

func (this *AIS) Parse(data string) error {
	this.last_measure.Raw = data

	var (
		split []string = strings.Split(data, ",")
		body  string   = split[5]
		runes []rune
		err   error
		value uint64
	)

	if len(split) < 7 {
		return errors.New("Message error")
	}

	//count_fragments
	value, err = strconv.ParseUint(split[1], 10, 8)
	if err != nil {
		return err
	}
	this.last_measure.Count_fragments = uint8(value)

	//fragment_number
	value, err = strconv.ParseUint(split[2], 10, 8)
	if err != nil {
		return err
	}
	this.last_measure.Fragment_number = uint8(value)

	//message_id
	if split[3] != "" {
		value, err = strconv.ParseUint(split[3], 10, 8)
		if err != nil {
			return err
		}
		this.last_measure.Message_id = uint8(value)
	}

	//radio channel
	this.Dynamic.Radio_channel = split[4]

	for i := 0; i < len(body); i++ {
		runes = append(runes, ascii_to_8bit(rune(body[i])))
	}

	this.last_measure.binary = dec_to_6bit(runes)

	err = this.parseBin()
	if err != nil {
		return err
	}

	err = this.checkCRC()

	return err
}

func (this AIS) GetMessageType() (uint8, string) {
	var description string

	switch this.last_measure.message_type {
	case 1:
		description = "Position Report Class A"
	case 2:
		description = "Position Report Class A (Assigned schedule)"
	case 3:
		description = "Position Report Class A (Response to interrogation)"
	case 4:
		description = "Base Station Report"
	case 5:
		description = "Ship and Voyage data"
	case 6:
		description = "Addressed Binary Message"
	case 7:
		description = "Binary Acknowledge"
	case 8:
		description = "Binary Broadcast Message"
	case 9:
		description = "Standard SAR Aircraft Position Report"
	case 10:
		description = "UTC and Date Inquiry"
	case 11:
		description = "UTC and Date Response"
	case 12:
		description = "Addressed Safety Related Message"
	case 13:
		description = "Safety Related Acknowledge"
	}

	return this.last_measure.message_type, description
}

func (this *AIS) parseBin() error {
	var (
		res []byte
		err error
	)

	// Message Type
	res = getBinary(this.last_measure.binary, 0, 6)
	this.last_measure.message_type = uint8(res[0])

	// Repeat Indicator
	res = getBinary(this.last_measure.binary, 6, 2)
	this.last_measure.repeat_indicator = uint8(res[0])

	// MMSI
	if this.Static.MMSI == 0 {
		res = getBinary(this.last_measure.binary, 8, 30)
		this.Static.MMSI = binary.BigEndian.Uint32(res)
	}

	switch this.last_measure.message_type {
	case 1, 2, 3:
		err = this.parsePositionReport()
	case 4:
		err = this.parseBaseStationReport()
	default:
		return errors.New("Undefined message type")
	}

	return err
}

func (this *AIS) checkCRC() error {
	var (
		body    string
		crc     int = 0
		index_s int = strings.Index(this.last_measure.Raw, "*")
	)

	if strings.Contains(this.last_measure.Raw, "!") && strings.Contains(this.last_measure.Raw, "*") {
		body = this.last_measure.Raw[1:index_s]
	} else {
		body = this.last_measure.Raw
	}

	for i := 0; i < len(body); i++ {
		if crc == 0 {
			crc = int(rune(body[i]))
		} else {
			crc ^= int(rune(body[i]))
		}
	}

	if strings.ToUpper(hex.EncodeToString([]byte{byte(crc)})) != strings.ToUpper(this.last_measure.Raw[index_s+1:index_s+3]) {
		return errors.New("CRC error")
	}

	return nil
}
