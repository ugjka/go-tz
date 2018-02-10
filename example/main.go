//Example program
package main

import (
	"fmt"

	"github.com/ugjka/go-tz"
)

func main() {
	// Example of loading custom GeoJSON
	// Shapefile Source: https://github.com/evansiroky/timezone-boundary-builder/

	/*f, err := os.Open("./dist/combined.json")
	if err != nil {
		panic(err)
	}
	err = gotz.LoadGeoJSON(f)
	if err != nil {
		panic(err)
	}
	*/

	// Loading Zone for Line Islands, Kiritimati
	zone, err := gotz.GetZone(gotz.Point{
		Lat: 1.74294, Lng: -157.21328,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(zone)
}
