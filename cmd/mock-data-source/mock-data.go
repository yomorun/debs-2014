/*
Copyright 2021 The YoMo Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"io"
	"log"
	"math/rand"
	"time"

	"github.com/yomorun/debs2014/internal/lib"
	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/client"
)

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

		data := []lib.Measurement{
			// plug 1
			{
				Id:          seconds,
				Timestamp:   seconds,
				Value:       rand.Float32() * 20,
				Property:    true, // load
				PlugId:      0,
				HouseholdId: 1,
				HouseId:     2,
			},
			{
				Id:          seconds,
				Timestamp:   seconds,
				Value:       rand.Float32() * 20,
				Property:    false, // work
				PlugId:      0,
				HouseholdId: 1,
				HouseId:     2,
			},

			// plug 2
			{
				Id:          seconds,
				Timestamp:   seconds,
				Value:       rand.Float32() * 20,
				Property:    true, // load
				PlugId:      3,
				HouseholdId: 1,
				HouseId:     2,
			},
			{
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
