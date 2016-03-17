package ais

import "time"

type Dynamic struct {
	Navigation_status Navigation_status
	Latitude          float64
	Longitude         float64
	Course            float32
	True_heading      uint16 // истинный курс. 0-359 градусов, 511 - недоступно
	Speed             float32
	Position_accuracy bool   // точность координат: 1 = high (<= 10 m), 0 = low (> 10 m) default
	Radio_channel     string // AIS Channel A is 161.975Mhz (87B); AIS Channel B is 162.025Mhz (88B)
	Time              time.Time
}
