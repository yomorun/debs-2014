package main

import (
	"context"
	"fmt"
	"log"
	"reflect"

	"github.com/yomorun/debs2014/internal/lib"

	"github.com/yomorun/yomo/pkg/client"
	"github.com/yomorun/yomo/pkg/rx"
)

var ss uint32 = 2
var t uint32 = 30

func main() {
	cli, err := client.NewServerless("load-prediction").Connect("localhost", 9000)
	if err != nil {
		log.Print("❌ Connect to zipper failure: ", err)
		return
	}

	defer cli.Close()
	cli.Pipe(Handler)
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
		x, ok := elem.(lib.Measurement)
		if !ok {
			err := fmt.Sprintf("expected type 'measurement', got '%v' instead",
				reflect.TypeOf(elem))
			fmt.Printf("[average] %v\n", err)
			return nil, fmt.Errorf(err)
		}

		if x.Property { // load
			plug := x.ToString()
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
			pred := (db[plug][idx] + lib.Median(data)) / 2
			fmt.Printf("[s_%v] %v %v\n", idx+2, plug, pred)
		}
	}
	fmt.Println("***************")

	idx += 1
	return 0.0, nil
}

// Query 1
func Handler(rxstream rx.RxStream) rx.RxStream {
	stream := rxstream.
		Subscribe(0x10).
		OnObserve(lib.Decoder).
		Map(lib.Printer).
		BufferWithTime(ss * 1e3).
		Map(average).
		Map(predict).
		Encode(0x1B)

	return stream
}
