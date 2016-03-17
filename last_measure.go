package ais

type last_measure struct {
	Raw              string
	Count_fragments  uint8 // количество фрагментов
	Fragment_number  uint8 // номер фрагмента
	Message_id       uint8 // идентификатор фрагментов
	message_type     uint8
	repeat_indicator uint8
	crc_valid        bool
	binary           []byte
}
