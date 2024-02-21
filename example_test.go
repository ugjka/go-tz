package tz_test

import (
	"fmt"

	"github.com/ugjka/go-tz/v2"
)

func ExampleGetZone() {
	// Loading Zone for Line Islands, Kiritimati
	p := tz.Point{
		Lon: -157.21328, Lat: 1.74294,
	}
	zone, err := tz.GetZone(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(zone[0])
	// Output: Pacific/Kiritimati
}
