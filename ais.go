package ais

import (
	//"bytes"
	//"encoding/binary"
	//"encoding/hex"
	//"errors"
	//"strconv"
	//"strings"
	"time"
)

const TIMEOUT_REMOVE_FRAGMENTS time.Duration = time.Second * 10

type AIS struct {
	MMSI              uint32
	IMO               uint32
	Message_type      uint8
	Radio_channel     string // AIS Channel A is 161.975Mhz (87B); AIS Channel B is 162.025Mhz (88B)
	Navigation_status uint8
	Speed             float32
	Position_accuracy bool // точность координат: 1 = high (<= 10 m), 0 = low (> 10 m) default
	Longitude         float64
	Latitude          float64
	Course            float32
	True_heading      uint16 // истинный курс. 0-359 градусов, 511 - недоступно
	Time              time.Time
	Type_of_EPFD      uint8
	Valid             bool
	is_parsed         bool
}

func (this AIS) IsParsed() bool {
	return this.is_parsed
}

//получить тип сообщения и его описание
func (this AIS) GetMessageType() string {
	var description string

	switch this.Message_type {
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

	return description
}

//получить тип сообщения и его описание
func (this AIS) GetNavigationStatus() string {
	var description string

	switch this.Navigation_status {
	case 0:
		description = "Under way using engine"
	case 1:
		description = "At anchor"
	case 2:
		description = "Not under command"
	case 3:
		description = "Restricted manoeuverability"
	case 4:
		description = "Constrained by her draught"
	case 5:
		description = "Moored"
	case 6:
		description = "Aground"
	case 7:
		description = "Engaged in Fishing"
	case 8:
		description = "Under way sailing"
	case 9:
		description = "Reserved for future amendment of Navigational Status for HSC"
	case 10:
		description = "Reserved for future amendment of Navigational Status for WIG"
	case 11, 12, 13, 14: // зарезервировано
		description = ""
	default:
		description = "Not defined"
	}

	return description
}

func (this AIS) GetTypeOfEPFD() string {
	var description string

	switch this.Type_of_EPFD {
	case 1:
		description = "GPS"
	case 2:
		description = "GLONASS"
	case 3:
		description = "Combined GPS/GLONASS"
	case 4:
		description = "Loran-C"
	case 5:
		description = "Chayka"
	case 6:
		description = "Integrated navigation system"
	case 7:
		description = "Surveyed"
	case 8:
		description = "Galileo"
	default:
		description = "Undefined"
	}

	return description
}

/*
func (this *AIS) Parse(data string) error {
	this.last_measure.Raw = data

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
		return errors.New("Message error")
	}

	//count_fragments
	value, err = strconv.ParseUint(split[1], 10, 8)
	if err != nil {
		return err
	}
	count_fragments = uint8(value)

	//fragment_number
	value, err = strconv.ParseUint(split[2], 10, 8)
	if err != nil {
		return err
	}
	fragment_number = uint8(value)

	//message_id
	if split[3] != "" {
		value, err = strconv.ParseUint(split[3], 10, 8)
		if err != nil {
			return err
		}
		message_id = uint8(value)
	}

	//radio channel
	this.Dynamic.Radio_channel = split[4]

	// если фрагментов несколько, считаем общую валидность CRC для всех фрагментов
	err = this.checkCRC()
	if count_fragments == 1 || fragment_number == 1 {
		this.crc_valid = err == nil
	} else {
		this.crc_valid = this.crc_valid && err == nil
	}

	for i := 0; i < len(body); i++ {
		runes = append(runes, ascii_to_8bit(rune(body[i])))
	}

	// TODO если фрагментов несколько, но пришли не все, то по таймауту удалить предыдущие

	if count_fragments > 1 && message_id != 0 {
		if len(this.fragments) == 0 {
			this.fragments = make(map[uint8][][]byte)
		}

		this.fragments[message_id] = append(this.fragments[message_id], dec_to_6bit(runes))

		// если пришли ещё не все фрагменты, пропустим парсинг
		if fragment_number == count_fragments {
			this.last_measure.binary = bytes.Join(this.fragments[message_id], nil)
			delete(this.fragments, message_id)
		} else {
			return nil
		}
	} else {
		this.last_measure.binary = dec_to_6bit(runes)
	}

	err = this.parseBin()
	if err != nil {
		return err
	}

	return err
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
	case 5:
		err = this.parseShipAndVoyage()
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
*/
