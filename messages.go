package ais

import (
	"encoding/binary"
	//"fmt"
	"time"
)

// 1, 2, 3
func (this *AIS) parsePositionReport() error {
	var res []byte
	var lat_lon uint32

	// Navigation Status
	res = getBinary(this.last_measure.binary, 38, 4)
	this.Dynamic.Navigation_status = GetNavigationStatus(uint8(res[0]))

	// Speed
	res = getBinary(this.last_measure.binary, 50, 10)
	this.Dynamic.Speed = float32(binary.BigEndian.Uint16(res)) / 10

	//Position_accuracy
	res = getBinary(this.last_measure.binary, 60, 1)
	this.Dynamic.Position_accuracy = res[0] == 0x01

	//Longitude
	res = getBinary(this.last_measure.binary, 61, 28)
	lat_lon = binary.BigEndian.Uint32(res) & 0x0FFFFFFF
	if lat_lon&0x08000000 > 0x00 {
		lat_lon = lat_lon ^ 0x0FFFFFFF
		this.Dynamic.Longitude = (float64(lat_lon+1) / 600000.0) * -1
	} else {
		this.Dynamic.Longitude = float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}

	//Latitude
	res = getBinary(this.last_measure.binary, 89, 27)
	lat_lon = binary.BigEndian.Uint32(res) & 0x07FFFFFF
	if lat_lon&0x04000000 > 0x00 {
		lat_lon = lat_lon ^ 0x07FFFFFF
		this.Dynamic.Latitude = (float64(lat_lon+1) / 600000.0) * -1
	} else {
		this.Dynamic.Latitude = float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}

	//Course
	res = getBinary(this.last_measure.binary, 116, 12)
	this.Dynamic.Course = float32(binary.BigEndian.Uint16(res)) / 10

	//True_heading
	res = getBinary(this.last_measure.binary, 128, 9)
	this.Dynamic.True_heading = uint16(binary.BigEndian.Uint16(res))

	return nil
}

// 4
func (this *AIS) parseBaseStationReport() error {
	var (
		res          []byte
		measure_time [6]int
		lat_lon      uint32
	)

	// Year
	res = getBinary(this.last_measure.binary, 38, 14)
	measure_time[0] = int(binary.BigEndian.Uint16(res))

	// Month
	res = getBinary(this.last_measure.binary, 52, 4)
	measure_time[1] = int(res[0])

	// Day
	res = getBinary(this.last_measure.binary, 56, 5)
	measure_time[2] = int(res[0])

	// Hour
	res = getBinary(this.last_measure.binary, 61, 5)
	measure_time[3] = int(res[0])

	// Minute
	res = getBinary(this.last_measure.binary, 66, 6)
	measure_time[4] = int(res[0])

	// Second
	res = getBinary(this.last_measure.binary, 72, 6)
	measure_time[5] = int(res[0])

	this.Dynamic.Time = time.Date(measure_time[0], time.Month(measure_time[1]), measure_time[2], measure_time[3], measure_time[4], measure_time[5], 0, time.UTC)

	//Longitude
	res = getBinary(this.last_measure.binary, 61, 28)
	lat_lon = binary.BigEndian.Uint32(res) & 0x0FFFFFFF
	if lat_lon&0x08000000 > 0x00 {
		lat_lon = lat_lon ^ 0x0FFFFFFF
		this.Dynamic.Longitude = (float64(lat_lon+1) / 600000.0) * -1
	} else {
		this.Dynamic.Longitude = float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}

	//Latitude
	res = getBinary(this.last_measure.binary, 89, 27)
	lat_lon = binary.BigEndian.Uint32(res) & 0x07FFFFFF
	if lat_lon&0x04000000 > 0x00 {
		lat_lon = lat_lon ^ 0x07FFFFFF
		this.Dynamic.Latitude = (float64(lat_lon+1) / 600000.0) * -1
	} else {
		this.Dynamic.Latitude = float64(binary.BigEndian.Uint32(res)) / 10000 / 60
	}

	return nil
}
