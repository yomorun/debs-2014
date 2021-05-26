package main

import (
	"context"
	"fmt"
	"log"
	"reflect"
	"sort"

	y3 "github.com/yomorun/y3-codec-golang"
	"github.com/yomorun/yomo/pkg/client"
	"github.com/yomorun/yomo/pkg/rx"
)

var ss uint32 = 2
var t uint32 = 30
var ws uint32 = 5

// Id           a unique identifier of the measurement
// Timestamp    timestamp of measurement (number of seconds since January 1, 1970, 00:00:00 GMT)
// Value        the measurement
// Property     type of the measurement: 0 for work or 1 for load
// PlugId      	a unique identifier (within a household) of the smart plug
// HouseholdId	a unique identifier of a household (within a house) where the plug is located
// HouseId     	a unique identifier of a house where the household with the plug is located
type Measurement struct {
	// TAGS ARE REQUIRED !!
	Id          uint32  `y3:"0x11"`
	Timestamp   uint32  `y3:"0x12"`
	Value       float32 `y3:"0x13"`
	Property    bool    `y3:"0x14"`
	PlugId      uint32  `y3:"0x15"`
	HouseholdId uint32  `y3:"0x16"`
	HouseId     uint32  `y3:"0x17"`
}

func (x Measurement) toString() string {
	// unique identifier
	return fmt.Sprintf("%v-%v-%v", x.PlugId, x.HouseholdId, x.HouseId)
}

func main() {
	cli, err := client.NewServerless("realtime-calc").Connect("localhost", 9000)
	if err != nil {
		log.Print("❌ Connect to zipper failure: ", err)
		return
	}

	defer cli.Close()
	cli.Pipe(Handler)
}

func decoder(v []byte) (interface{}, error) {
	var mold Measurement

	// defined in y3-codec-golang/types.go
	// decode []byte to interface{}
	err := y3.ToObject(v, &mold)
	if err != nil {
		fmt.Printf("[decoder] %v\n", err)
		return nil, err
	}
	return mold, nil
}

// an empty interface (interface{}) may hold values of any type;
// empty interfaces are used by code that handles values of unknown type
func printer(_ context.Context, i interface{}) (interface{}, error) {
	x, ok := i.(Measurement)
	if !ok {
		err := fmt.Sprintf("expected type 'measurement', got '%v' instead",
			reflect.TypeOf(i))
		fmt.Printf("[printer] %v\n", err)
		return nil, fmt.Errorf(err)
	}

	var prop string
	if x.Property {
		prop = "load"
	} else {
		prop = "work"
	}
	fmt.Printf("[%v] %v %v %v\n",
		x.Timestamp, x.Value, x.toString(), prop)
	return i, nil
}

// plug # -> slice # -> avg
var db = make(map[string]map[uint32]float32)
var idx uint32 = 0

func average(_ context.Context, i interface{}) (interface{}, error) {
	lst, ok := i.([]interface{})
	if !ok {
		err := fmt.Sprintf("expected type '[]interface{}', got '%v' instead",
			reflect.TypeOf(i))
		fmt.Printf("[average] %v\n", err)
		return nil, fmt.Errorf(err)
	}

	// plug # -> value
	total := make(map[string]float32)
	count := make(map[string]float32)

	for _, elem := range lst {
		x, ok := elem.(Measurement)
		if !ok {
			err := fmt.Sprintf("expected type 'measurement', got '%v' instead",
				reflect.TypeOf(elem))
			fmt.Printf("[average] %v\n", err)
			return nil, fmt.Errorf(err)
		}

		if x.Property { // load
			plug := x.toString()
			total[plug] += x.Value
			count[plug] += 1.0
		}
	}

	// save to db
	fmt.Println("*** average ***")
	for plug, v := range total {
		avg := v / count[plug]
		fmt.Printf("[s_%v] %v %v\n", idx, plug, avg)

		_, ok := db[plug]
		if !ok {
			db[plug] = make(map[uint32]float32)
		}
		db[plug][idx] = avg
	}
	fmt.Println("***************")
	return i, nil
}

func predict(_ context.Context, i interface{}) (interface{}, error) {
	// j = i + 2 - n * k
	// 	   k is the number of slices in a 24 hour period [tentative]
	//     n is a natural number with values between 1 and floor((i + 2) / k)
	k := t / ss

	fmt.Println("*** predict ***")
	l := (idx + 2) / k
	if l == 0 {
		fmt.Println("not enough data")
	} else {
		// possible values for j
		lst := make([]uint32, l)
		for m := range lst {
			n := uint32(m + 1)
			j := idx + 2 - n*k
			lst[m] = j
		}

		for plug := range db {
			// average load for s_j
			data := make([]float32, l)
			for m, j := range lst {
				data[m] = db[plug][j]
			}
			pred := (db[plug][idx] + median(data)) / 2
			fmt.Printf("[s_%v] %v %v\n", idx+2, plug, pred)
		}
	}
	fmt.Println("***************")

	idx += 1
	return 0.0, nil
}

func outliers(_ context.Context, i interface{}) (interface{}, error) {
	lst, ok := i.([]interface{})
	if !ok {
		err := fmt.Sprintf("expected type '[]interface{}', got '%v' instead",
			reflect.TypeOf(i))
		fmt.Printf("[outliers] %v\n", err)
		return nil, fmt.Errorf(err)
	}

	all := make([]float32, 0, len(lst))
	indiv := make(map[string][]float32) // plug # -> values

	for _, elem := range lst {
		x, ok := elem.(Measurement)
		if !ok {
			err := fmt.Sprintf("expected type 'measurement', got '%v' instead",
				reflect.TypeOf(elem))
			fmt.Printf("[outliers] %v\n", err)
			return nil, fmt.Errorf(err)
		}

		if x.Property { // load
			all = append(all, x.Value)

			plug := x.toString()
			indiv[plug] = append(indiv[plug], x.Value)
		}
	}

	v := median(all)
	fmt.Printf("all plugs: %v\n", v)

	fmt.Println("*** outliers ***")
	for plug, vs := range indiv {
		m := median(vs)
		if median(vs) > v {
			fmt.Printf("[w_%v] %v %v\n", idx, plug, m)
		}
	}
	fmt.Println("****************")

	idx += 1
	return 0.0, nil
}

// pre: len(data) > 0
func median(data []float32) float32 {
	l := len(data)
	sort.Slice(data, func(i, j int) bool {
		return data[i] < data[j]
	})

	if l%2 == 0 {
		return (data[l/2-1] + data[l/2]) / 2
	} else {
		return data[l/2]
	}
}

// // Query 1
// func Handler(rxstream rx.RxStream) rx.RxStream {
// 	stream := rxstream.
// 		Subscribe(0x10).
// 		OnObserve(decoder).
// 		Map(printer).
// 		BufferWithTime(ss * 1e3).
// 		Map(average).
// 		Map(predict).
// 		Encode(0x11)

// 	return stream
// }

// Query 2
func Handler(rxstream rx.RxStream) rx.RxStream {
	stream := rxstream.
		Subscribe(0x10).
		OnObserve(decoder).
		Map(printer).
		BufferWithTime(ws * 1e3).
		Map(outliers).
		Encode(0x11)

	return stream
}
