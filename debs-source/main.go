package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/client"
)

// Id           a unique identifier of the measurement
// Timestamp    timestamp of measurement (number of seconds since January 1, 1970, 00:00:00 GMT)
// Value        the measurement
// Property     type of the measurement: 0 for work or 1 for load
// PlugId          a unique identifier (within a household) of the smart plug
// HouseholdId    a unique identifier of a household (within a house) where the plug is located
// HouseId         a unique identifier of a house where the household with the plug is located
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
		log.Println(err)
		return
	}
	log.Println("connected")
	generateData(client)
}

var codec = y3.NewCodec(0x10)

func parseUint32(s string) uint32 {
	v, err := strconv.ParseUint(s, 10, 32)
	if err != nil {
		log.Println("unable to parse string as uint32")
		return 0
	}
	return uint32(v)
}

func parseFloat32(s string) float32 {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		log.Println("unable to parse string as float32")
		return 0
	}
	return float32(v)
}

func parseBool(s string) bool {
	return s == "1"
}

func generateData(stream io.Writer) {
	log.Println("reading data file...")
	filename := "data.csv"
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	lines, err := csv.NewReader(file).ReadAll()
	if err != nil {
		fmt.Println(err)
	}

	data := make(map[uint32][]measurement)
	// keys := make([]uint32, 0)
	for _, line := range lines {
		x := measurement{
			Id:          parseUint32(line[0]),
			Timestamp:   parseUint32(line[1]),
			Value:       parseFloat32(line[2]),
			Property:    parseBool(line[3]),
			PlugId:      parseUint32(line[4]),
			HouseholdId: parseUint32(line[5]),
			HouseId:     parseUint32(line[6]),
		}
		data[x.Timestamp] = append(data[x.Timestamp], x)
		// keys = append(keys, x.Timestamp)
	}

	keys := make([]uint32, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	log.Println("sending data...")
	for _, ts := range keys {
		start := time.Now()
		log.Println(ts)
		arr := data[ts]

		for _, x := range arr {
			buf, _ := codec.Marshal(x)
			_, err := stream.Write(buf)
			if err != nil {
				log.Println(err)
			} else {
				log.Printf("%v-%v-%v: %v",
					x.PlugId, x.HouseholdId, x.HouseId, x.Value)
			}
		}
		t := time.Now()
		elapsed := t.Sub(start)
		time.Sleep(time.Second - elapsed)
	}
	log.Println("done")
}
