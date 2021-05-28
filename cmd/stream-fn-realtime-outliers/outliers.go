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

var ws uint32 = 5

func main() {
	cli, err := client.NewServerless("outliers").Connect("localhost", 9000)
	if err != nil {
		log.Print("❌ Connect to zipper failure: ", err)
		return
	}

	defer cli.Close()
	cli.Pipe(Handler)
}

// an empty interface (interface{}) may hold values of any type;
// empty interfaces are used by code that handles values of unknown type
func printer(_ context.Context, i interface{}) (interface{}, error) {
	x, ok := i.(lib.Measurement)
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
		x.Timestamp, x.Value, x.ToString(), prop)
	return i, nil
}

// plug # -> slice # -> avg
var idx uint32 = 0

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
		x, ok := elem.(lib.Measurement)
		if !ok {
			err := fmt.Sprintf("expected type 'measurement', got '%v' instead",
				reflect.TypeOf(elem))
			fmt.Printf("[outliers] %v\n", err)
			return nil, fmt.Errorf(err)
		}

		if x.Property { // load
			all = append(all, x.Value)

			plug := x.ToString()
			indiv[plug] = append(indiv[plug], x.Value)
		}
	}

	v := lib.Median(all)
	fmt.Printf("all plugs: %v\n", v)

	fmt.Println("*** outliers ***")
	for plug, vs := range indiv {
		m := lib.Median(vs)
		if lib.Median(vs) > v {
			fmt.Printf("[w_%v] %v %v\n", idx, plug, m)
		}
	}
	fmt.Println("****************")

	idx += 1
	return 0.0, nil
}

// Query 2
func Handler(rxstream rx.RxStream) rx.RxStream {
	stream := rxstream.
		Subscribe(0x10).
		OnObserve(lib.Decoder).
		Map(printer).
		BufferWithTime(ws * 1e3).
		Map(outliers).
		Encode(0x1A)

	return stream
}
