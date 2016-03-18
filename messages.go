package ais

import (
	"encoding/binary"
	//"fmt"
	"time"
)

// 1, 2, 3
func (this Parser) parsePositionReport(ais_item *AIS, bits []byte) {
	var res []byte
	var lat_lon uint32

	// Navigation Status
	res = getBinary(bits, 38, 4)
	ais_item.Navigation_status = uint8(res[0])

	// Speed
	res = getBinary(bits, 50, 10)
	ais_item.Speed = float32(binary.BigEndian.Uint16(res)) / 10

	//Position_accuracy
	res = getBinary(bits, 60, 1)
	ais_item.Position_accuracy = res[0] == 0x01

	//Longitude
	res = getBinary(bits, 61, 28)
	lat_lon = binary.BigEndian.Uint32(res) & 0x0FFFFFFF
	if lat_lon&0x08000000 > 0x00 {
		lat_lon = lat_lon ^ 0x0FFFFFFF
		ais_item.Longitude = (float64(lat_lon+1) / 600000.0) * -1
	} else {
		ais_item.Longitude = float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}

	//Latitude
	res = getBinary(bits, 89, 27)
	lat_lon = binary.BigEndian.Uint32(res) & 0x07FFFFFF
	if lat_lon&0x04000000 > 0x00 {
		lat_lon = lat_lon ^ 0x07FFFFFF
		ais_item.Latitude = (float64(lat_lon+1) / 600000.0) * -1
	} else {
		ais_item.Latitude = float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}

	//Course
	res = getBinary(bits, 116, 12)
	ais_item.Course = float32(binary.BigEndian.Uint16(res)) / 10

	//True_heading
	res = getBinary(bits, 128, 9)
	ais_item.True_heading = uint16(binary.BigEndian.Uint16(res))
}

// 4
func (this Parser) parseBaseStationReport(ais_item *AIS, bits []byte) {
	var (
		res          []byte
		measure_time [6]int
		lat_lon      uint32
	)

	// Year
	res = getBinary(bits, 38, 14)
	measure_time[0] = int(binary.BigEndian.Uint16(res))

	// Month
	res = getBinary(bits, 52, 4)
	measure_time[1] = int(res[0])

	// Day
	res = getBinary(bits, 56, 5)
	measure_time[2] = int(res[0])

	// Hour
	res = getBinary(bits, 61, 5)
	measure_time[3] = int(res[0])

	// Minute
	res = getBinary(bits, 66, 6)
	measure_time[4] = int(res[0])

	// Second
	res = getBinary(bits, 72, 6)
	measure_time[5] = int(res[0])

	ais_item.Time = time.Date(measure_time[0], time.Month(measure_time[1]), measure_time[2], measure_time[3], measure_time[4], measure_time[5], 0, time.UTC)

	//Longitude
	res = getBinary(bits, 61, 28)
	lat_lon = binary.BigEndian.Uint32(res) & 0x0FFFFFFF
	if lat_lon&0x08000000 > 0x00 {
		lat_lon = lat_lon ^ 0x0FFFFFFF
		ais_item.Longitude = (float64(lat_lon+1) / 600000.0) * -1
	} else {
		ais_item.Longitude = float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}

	//Latitude
	res = getBinary(bits, 89, 27)
	lat_lon = binary.BigEndian.Uint32(res) & 0x07FFFFFF
	if lat_lon&0x04000000 > 0x00 {
		lat_lon = lat_lon ^ 0x07FFFFFF
		ais_item.Latitude = (float64(lat_lon+1) / 600000.0) * -1
	} else {
		ais_item.Latitude = float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}

	res = getBinary(bits, 134, 4)
	ais_item.Type_of_EPFD = uint8(res[0])
}

// 5
func (this Parser) parseShipAndVoyage(ais_item *AIS, bits []byte) {
	var res []byte

	// IMO
	res = getBinary(bits, 40, 30)
	ais_item.IMO = binary.BigEndian.Uint32(res)

	//fmt.Println(res)
}
