package main

import (
	"fmt"

	"github.com/ugjka/go-tz"
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
