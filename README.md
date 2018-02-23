# go-tz

tz-lookup by lng and lat

[![GoDoc](https://godoc.org/github.com/ugjka/go-tz?status.svg)](https://godoc.org/github.com/ugjka/go-tz)
[![Go Report Card](https://goreportcard.com/badge/github.com/ugjka/go-tz)](https://goreportcard.com/report/github.com/ugjka/go-tz)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)


lookup timezone for a given location

```go
// Loading Zone for Line Islands, Kiritimati
	zone, err := gotz.GetZone(gotz.Point{
		Lat: 1.74294, Lon: -157.21328,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(zone)
  ```
  
  ```
  [ugjka@archee example]$ go run main.go 
Pacific/Kiritimati
```

Uses simplified shapefile from https://github.com/evansiroky/timezone-boundary-builder/

GeoJson Simplification done with http://mapshaper.org/ and https://github.com/foursquare/shapefile-geo
