package ais

import (
	"encoding/binary"
	"time"
)

func (this Parser) getLon(res []byte) float64 {
	lat_lon := binary.BigEndian.Uint32(res) & 0x0FFFFFFF
	if lat_lon&0x08000000 > 0x00 {
		lat_lon = lat_lon ^ 0x0FFFFFFF
		return (float64(lat_lon+1) / 600000.0) * -1
	} else {
		return float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}
}

func (this Parser) getLat(res []byte) float64 {
	lat_lon := binary.BigEndian.Uint32(res) & 0x07FFFFFF
	if lat_lon&0x04000000 > 0x00 {
		lat_lon = lat_lon ^ 0x07FFFFFF
		return (float64(lat_lon+1) / 600000.0) * -1
	} else {
		return float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}
}

// 1, 2, 3
func (this Parser) parsePositionReport(ais_item *AIS, bits Bits_array) {
	var res []byte

	// Navigation Status
	res = bits.GetBinary(38, 4)
	ais_item.Navigation_status = uint8(res[0])

	// Speed
	res = bits.GetBinary(50, 10)
	ais_item.Speed = float32(binary.BigEndian.Uint16(res)) / 10

	//Position_accuracy
	res = bits.GetBinary(60, 1)
	ais_item.Position_accuracy = res[0] == 0x01

	//Longitude
	res = bits.GetBinary(61, 28)
	ais_item.Longitude = this.getLon(res)

	//Latitude
	res = bits.GetBinary(89, 27)
	ais_item.Latitude = this.getLat(res)

	//Course
	res = bits.GetBinary(116, 12)
	ais_item.Course = float32(binary.BigEndian.Uint16(res)) / 10

	//True_heading
	res = bits.GetBinary(128, 9)
	ais_item.True_heading = binary.BigEndian.Uint16(res)
}

// 4
func (this Parser) parseBaseStationReport(ais_item *AIS, bits Bits_array) {
	var (
		res          []byte
		measure_time [6]int
	)

	// Year
	res = bits.GetBinary(38, 14)
	measure_time[0] = int(binary.BigEndian.Uint16(res))

	// Month
	res = bits.GetBinary(52, 4)
	measure_time[1] = int(res[0])

	// Day
	res = bits.GetBinary(56, 5)
	measure_time[2] = int(res[0])

	// Hour
	res = bits.GetBinary(61, 5)
	measure_time[3] = int(res[0])

	// Minute
	res = bits.GetBinary(66, 6)
	measure_time[4] = int(res[0])

	// Second
	res = bits.GetBinary(72, 6)
	measure_time[5] = int(res[0])

	ais_item.Time = time.Date(measure_time[0], time.Month(measure_time[1]), measure_time[2], measure_time[3], measure_time[4], measure_time[5], 0, time.UTC)

	//Longitude
	res = bits.GetBinary(79, 28)
	ais_item.Longitude = this.getLon(res)

	//Latitude
	res = bits.GetBinary(107, 27)
	ais_item.Latitude = this.getLat(res)

	res = bits.GetBinary(134, 4)
	ais_item.Type_of_EPFD = uint8(res[0])
}

// 5
func (this Parser) parseShipAndVoyage(ais_item *AIS, bits Bits_array) {
	var res []byte

	// IMO
	res = bits.GetBinary(40, 30)
	ais_item.IMO = binary.BigEndian.Uint32(res)

	//Call Sign
	ais_item.Call_sign = bits.GetText(70, 42)

	//Ship_name
	ais_item.Ship_name = bits.GetText(112, 120)

	// Ship Type
	res = bits.GetBinary(232, 8)
	ais_item.Ship_type = uint8(res[0])

	// Dimension_to_bow
	res = bits.GetBinary(240, 9)
	ais_item.Dimension_to_bow = binary.BigEndian.Uint16(res)

	// Dimension_to_stern
	res = bits.GetBinary(249, 9)
	ais_item.Dimension_to_stern = binary.BigEndian.Uint16(res)

	// Dimension_to_port
	res = bits.GetBinary(258, 6)
	ais_item.Dimension_to_port = uint8(res[0])

	// Dimension_to_starboard
	res = bits.GetBinary(264, 6)
	ais_item.Dimension_to_starboard = uint8(res[0])

	res = bits.GetBinary(270, 4)
	ais_item.Type_of_EPFD = uint8(res[0])

	// ETA month
	res = bits.GetBinary(274, 4)
	ais_item.ETA[0] = uint8(res[0]) // 0 = N/A

	// ETA day
	res = bits.GetBinary(278, 5)
	ais_item.ETA[1] = uint8(res[0]) // 0 = N/A

	// ETA hour
	res = bits.GetBinary(283, 5)
	ais_item.ETA[2] = uint8(res[0]) // 24 = N/A

	// ETA minute
	res = bits.GetBinary(288, 6) // 60 = N/A
	ais_item.ETA[3] = uint8(res[0])

	// Draught
	res = bits.GetBinary(294, 8)
	ais_item.Draught = float32(res[0]) / 10

	//Destination
	ais_item.Destination = bits.GetText(302, 120)
}

// 18
func (this Parser) parseStandarBPositionReport(ais_item *AIS, bits Bits_array) {
	var res []byte

	// Speed
	res = bits.GetBinary(46, 10)
	ais_item.Speed = float32(binary.BigEndian.Uint16(res)) / 10

	//Position_accuracy
	res = bits.GetBinary(56, 1)
	ais_item.Position_accuracy = res[0] == 0x01

	//Longitude
	res = bits.GetBinary(57, 28)
	ais_item.Longitude = this.getLon(res)

	//Latitude
	res = bits.GetBinary(85, 27)
	ais_item.Latitude = this.getLat(res)

	//Course
	res = bits.GetBinary(112, 12)
	ais_item.Course = float32(binary.BigEndian.Uint16(res)) / 10

	//True_heading
	res = bits.GetBinary(124, 9)
	ais_item.True_heading = binary.BigEndian.Uint16(res)
}

// 19
func (this Parser) parseExtendedBPositionReport(ais_item *AIS, bits Bits_array) {
	var res []byte

	// Speed
	res = bits.GetBinary(46, 10)
	ais_item.Speed = float32(binary.BigEndian.Uint16(res)) / 10

	//Position_accuracy
	res = bits.GetBinary(56, 1)
	ais_item.Position_accuracy = res[0] == 0x01

	//Longitude
	res = bits.GetBinary(57, 28)
	ais_item.Longitude = this.getLon(res)

	//Latitude
	res = bits.GetBinary(85, 27)
	ais_item.Latitude = this.getLat(res)

	//Course
	res = bits.GetBinary(112, 12)
	ais_item.Course = float32(binary.BigEndian.Uint16(res)) / 10

	//True_heading
	res = bits.GetBinary(124, 9)
	ais_item.True_heading = binary.BigEndian.Uint16(res)

	//Ship_name
	ais_item.Ship_name = bits.GetText(143, 120)

	// Ship Type
	res = bits.GetBinary(263, 8)
	ais_item.Ship_type = uint8(res[0])

	// Dimension_to_bow
	res = bits.GetBinary(271, 9)
	ais_item.Dimension_to_bow = binary.BigEndian.Uint16(res)

	// Dimension_to_stern
	res = bits.GetBinary(280, 9)
	ais_item.Dimension_to_stern = binary.BigEndian.Uint16(res)

	// Dimension_to_port
	res = bits.GetBinary(289, 6)
	ais_item.Dimension_to_port = uint8(res[0])

	// Dimension_to_starboard
	res = bits.GetBinary(295, 6)
	ais_item.Dimension_to_starboard = uint8(res[0])

	res = bits.GetBinary(301, 4)
	ais_item.Type_of_EPFD = uint8(res[0])
}

// 24
func (this Parser) parseC_CSStaticDataReport(ais_item *AIS, bits Bits_array) {
	var res []byte

	//Ship_name
	ais_item.Ship_name = bits.GetText(40, 120)

	// Ship Type
	res = bits.GetBinary(40, 8)
	ais_item.Ship_type = uint8(res[0])

	//Call_sign
	ais_item.Call_sign = bits.GetText(90, 42)

	// Dimension_to_bow
	res = bits.GetBinary(132, 9)
	ais_item.Dimension_to_bow = binary.BigEndian.Uint16(res)

	// Dimension_to_stern
	res = bits.GetBinary(141, 9)
	ais_item.Dimension_to_stern = binary.BigEndian.Uint16(res)

	// Dimension_to_port
	res = bits.GetBinary(150, 6)
	ais_item.Dimension_to_port = uint8(res[0])

	// Dimension_to_starboard
	res = bits.GetBinary(156, 6)
	ais_item.Dimension_to_starboard = uint8(res[0])
}
