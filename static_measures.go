package ais

type Static struct {
	MMSI         uint32
	IMO          uint32
	Type_of_EPFD uint8
}

func (this Static) GetTypeOfEPFD() (uint8, string) {
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

	return this.Type_of_EPFD, description
}
