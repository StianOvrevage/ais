package ais

import (
	"time"
)

// TODO
const TIMEOUT_REMOVE_FRAGMENTS time.Duration = time.Second * 10

type AIS struct {
	MMSI                   uint32
	IMO                    uint32
	Message_type           uint8
	Radio_channel          string // AIS Channel A is 161.975Mhz (87B); AIS Channel B is 162.025Mhz (88B)
	Call_sign              string // позывной
	Ship_name              string // название судна
	Ship_type              uint8  // тип судна
	Navigation_status      uint8
	Speed                  float32
	Position_accuracy      bool // точность координат: 1 = high (<= 10 m), 0 = low (> 10 m) default
	Longitude              float64
	Latitude               float64
	Course                 float32 // вектор направления
	True_heading           uint16  // истинный курс. 0-359 градусов, 511 - недоступно
	Draught                float32 // осадка судна, м
	Destination            string
	Time                   time.Time
	ETA                    [4]uint8 // ожидаемое время прибытия [month, day, hour, minute]
	Type_of_EPFD           uint8
	Dimension_to_bow       uint16 // Dimension_to_bow + Dimension_to_stern = длинна судна
	Dimension_to_stern     uint16
	Dimension_to_port      uint8 // Dimension_to_port + Dimension_to_starboard = ширина судна
	Dimension_to_starboard uint8
	Valid                  bool
	is_parsed              bool
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

func (this AIS) GetShipType() string {
	var description string

	switch this.Ship_type {
	case 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19: //Reserved for future use
		description = ""
	case 20:
		description = "Wing in ground (WIG), all ships ofthis type"
	case 21:
		description = "Wing in ground (WIG), Hazardous category A"
	case 22:
		description = "Wing in ground (WIG), Hazardous category B"
	case 23:
		description = "Wing in ground (WIG), Hazardous category C"
	case 24:
		description = "Wing in ground (WIG), Hazardous category D"
	case 30:
		description = "Fishing"
	case 31:
		description = "Towing"
	case 32:
		description = "Towing: length exceeds 200m or breadth exceeds 25m"
	case 33:
		description = "Dredging or underwater ops"
	case 34:
		description = "Diving ops"
	case 35:
		description = "Military ops"
	case 36:
		description = "Sailing"
	case 37:
		description = "Pleasure Craft"
	case 38, 39: // Reserved
		description = ""
	case 40:
		description = "High speed craft (HSC), all ships of this type"
	case 41:
		description = "High speed craft (HSC), Hazardous category A"
	case 42:
		description = "High speed craft (HSC), Hazardous category B"
	case 43:
		description = "High speed craft (HSC), Hazardous category C"
	case 44:
		description = "High speed craft (HSC), Hazardous category D"
	case 45, 46, 47, 48: // Reserved
		description = "High speed craft (HSC)"
	case 49:
		description = "High speed craft (HSC), No additional information"
	case 50:
		description = "Pilot Vessel"
	case 51:
		description = "Search and Rescue vessel"
	case 52:
		description = "Tug"
	case 53:
		description = "Port Tender"
	case 54:
		description = "Anti-pollution equipment"
	case 55:
		description = "Law Enforcement"
	case 56, 57:
		description = "Spare - Local Vessel"
	case 58:
		description = "Medical Transport"
	case 59:
		description = "Ship according to RR Resolution No.18"
	case 60:
		description = "Passenger, all ships of this type"
	case 61:
		description = "Passenger, Hazardous category A"
	case 62:
		description = "Passenger, Hazardous category B"
	case 63:
		description = "Passenger, Hazardous category C"
	case 64:
		description = "Passenger, Hazardous category D"
	case 65, 66, 67, 68: // Reserved
		description = "Passenger"
	case 69:
		description = "Passenger, No additional information"
	case 70:
		description = "Cargo, all ships of this type"
	case 71:
		description = "Cargo, Hazardous category A"
	case 72:
		description = "Cargo, Hazardous category B"
	case 73:
		description = "Cargo, Hazardous category C"
	case 74:
		description = "Cargo, Hazardous category D"
	case 75, 76, 77, 78: // Reserved
		description = "Cargo"
	case 79:
		description = "Cargo, No additional information"
	case 80:
		description = "Tanker, all ships of this type"
	case 81:
		description = "Tanker, Hazardous category A"
	case 82:
		description = "Tanker, Hazardous category B"
	case 83:
		description = "Tanker, Hazardous category C"
	case 84:
		description = "Tanker, Hazardous category D"
	case 85, 86, 87, 88: // Reserved
		description = "Tanker"
	case 89:
		description = "Tanker, No additional information"
	case 90:
		description = "Other Type, all ships of this type"
	case 91:
		description = "Other Type, Hazardous category A"
	case 92:
		description = "Other Type, Hazardous category B"
	case 93:
		description = "Other Type, Hazardous category C"
	case 94:
		description = "Other Type, Hazardous category D"
	case 95, 96, 97, 98: // Reserved
		description = "Other Type"
	case 99:
		description = "Other Type, no additional information"
	default:
		description = "Not available"
	}

	return description
}
