package ais

type Navigation_status struct {
	Id          uint8
	Description string
}

func GetNavigationStatus(id uint8) Navigation_status {
	var description string

	switch id {
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

	var status Navigation_status = Navigation_status{Id: id, Description: description}

	return status
}
