package main

import (
	"io"
	"log"
	"math/rand"
	"time"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/client"
)

// Id           a unique identifier of the measurement
// Timestamp    timestamp of measurement (number of seconds since January 1, 1970, 00:00:00 GMT)
// Value        the measurement
// Property     type of the measurement: 0 for work or 1 for load
// PlugId      	a unique identifier (within a household) of the smart plug
// HouseholdId	a unique identifier of a household (within a house) where the plug is located
// HouseId     	a unique identifier of a house where the household with the plug is located
type measurement struct {
	Id          uint32  `y3:"0x11"`
	Timestamp   uint32  `y3:"0x12"`
	Value       float32 `y3:"0x13"`
	Property    bool    `y3:"0x14"`
	PlugId      uint32  `y3:"0x15"`
	HouseholdId uint32  `y3:"0x16"`
	HouseId     uint32  `y3:"0x17"`
}

func main() {
	client, err := client.NewSource("debs-source").Connect("localhost", 9000)
	if err != nil {
		log.Print(err)
		return
	}
	log.Print("connected")
	generateData(client)
}

var codec = y3.NewCodec(0x10)

func generateData(stream io.Writer) {
	log.Print("generating data...")
	for {
		start := time.Now()
		seconds := uint32(start.Unix())
		// millis := start.UnixNano() / 1e6

		data := []measurement{
			// plug 1
			measurement{
				Id:          seconds,
				Timestamp:   seconds,
				Value:       rand.Float32() * 20,
				Property:    true, // load
				PlugId:      0,
				HouseholdId: 1,
				HouseId:     2,
			},
			measurement{
				Id:          seconds,
				Timestamp:   seconds,
				Value:       rand.Float32() * 20,
				Property:    false, // work
				PlugId:      0,
				HouseholdId: 1,
				HouseId:     2,
			},

			// plug 2
			measurement{
				Id:          seconds,
				Timestamp:   seconds,
				Value:       rand.Float32() * 20,
				Property:    true, // load
				PlugId:      3,
				HouseholdId: 1,
				HouseId:     2,
			},
			measurement{
				Id:          seconds,
				Timestamp:   seconds,
				Value:       rand.Float32() * 20,
				Property:    false, //work
				PlugId:      3,
				HouseholdId: 1,
				HouseId:     2,
			},
		}

		for _, x := range data {
			buf, _ := codec.Marshal(x)
			_, err := stream.Write(buf)
			if err != nil {
				log.Print(err)
			} else {
				log.Printf("[%v] %v %v-%v-%v", x.Id, x.Value,
					x.PlugId, x.HouseholdId, x.HouseId)
			}
		}
		t := time.Now()
		elapsed := t.Sub(start)
		time.Sleep(time.Second - elapsed)
	}
}
