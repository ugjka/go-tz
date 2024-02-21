package main

import (
	"fmt"

	"gopkg.in/ugjka/go-tz.v2"
)

func main() {
	zone, err := tz.GetZone(tz.Point{
		Lon: -157.21328, Lat: 1.74294,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(zone[0])
}
