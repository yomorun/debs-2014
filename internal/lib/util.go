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

package lib

import (
	"context"
	"fmt"
	"reflect"
	"sort"

	"github.com/yomorun/y3-codec-golang"
)

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

func (x Measurement) ToString() string {
	// unique identifier
	return fmt.Sprintf("%v-%v-%v", x.PlugId, x.HouseholdId, x.HouseId)
}

// Deserialize data from stream
func Decoder(v []byte) (interface{}, error) {
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
func Printer(_ context.Context, i interface{}) (interface{}, error) {
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
		x.Timestamp, x.Value, x.ToString(), prop)
	return i, nil
}

// pre: len(data) > 0
func Median(data []float32) float32 {
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
