# Golang ais

Парсинг сообщений Automatic Identification System

> go get github.com/mil-ast/ais

### Пример
```go
package main

import(
	"fmt"
	"github.com/mil-ast/ais"
)

func main() {
	var item_ais ais.AIS
	var end bool
	
	var messages []string = []string{
		"!AIVDM,1,1,,A,15MgK45P3@G?fl0E`JbR0OwT0@MS,0*4E", // 1
		"!AIVDM,1,1,,A,402PeI1uho;Q=OL9>LE>dF1000S:,0*01", // 4
		"!AIVDM,2,1,4,A,5815<9h1aLU1KMPs800<iD:0P58ltqT00000000t3jBA<4nE0DjT83lh,0*23", // 5
		"!AIVDM,2,2,4,A,H<U000000000000,2*31", // 5
	}
	
	parser := ais.GetParser()
	
	var (
		description string
	)
	
	for i := 0; i < len(messages); i++ {
		item_ais, end = parser.Parse(messages[i])
		if !end {
			continue
		}
		
		description = item_ais.GetMessageType()

		println(fmt.Sprintf("Message type %d (%s)", item_ais.Message_type, description))
		println(fmt.Sprintf("MMSI (%d)", item_ais.MMSI))
		println(fmt.Sprintf("IMO (%d)", item_ais.IMO))
		println(fmt.Sprintf("Radio_channel (%s)", item_ais.Radio_channel))
		
		description = item_ais.GetNavigationStatus()
		
		println(fmt.Sprintf("Navigation_status %d (%s)", item_ais.Navigation_status, description))
		println(fmt.Sprintf("Speed (%f)", item_ais.Speed))
		println(fmt.Sprintf("Position_accuracy (%v)", item_ais.Position_accuracy))
		println(fmt.Sprintf("Latitude (%f)", item_ais.Latitude))
		println(fmt.Sprintf("Longitude (%f)", item_ais.Longitude))
		println(fmt.Sprintf("Course (%f)", item_ais.Course))
		println(fmt.Sprintf("True_heading (%d)", item_ais.True_heading))
		println(fmt.Sprintf("Time (%v)", item_ais.Time))
		
		description = item_ais.GetTypeOfEPFD()
		
		println(fmt.Sprintf("Type_of_EPFD %d (%s)", item_ais.Type_of_EPFD, description))
		
		println("-----------------------------")
	}
}
```