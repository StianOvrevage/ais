package ais

import (
	"time"
)

// TODO
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
