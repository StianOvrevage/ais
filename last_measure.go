package ais

type last_measure struct {
	Raw              string
	message_type     uint8
	repeat_indicator uint8
	binary           []byte
}
