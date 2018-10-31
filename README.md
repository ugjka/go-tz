# go-tz

tz-lookup by lon and lat

[![GoDoc](https://godoc.org/github.com/ugjka/go-tz?status.svg)](https://godoc.org/github.com/ugjka/go-tz)
[![Go Report Card](https://goreportcard.com/badge/github.com/ugjka/go-tz)](https://goreportcard.com/report/github.com/ugjka/go-tz)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Donate](https://dl.ugjka.net/Donate-PayPal-green.svg)](https://www.paypal.com/cgi-bin/webscr?cmd=_s-xclick&hosted_button_id=UVTCZYQ3FVNCY)

lookup timezone for a given location

```go
// Loading Zone for Line Islands, Kiritimati
zone, err := gotz.GetZone(gotz.Point{
    Lat: 1.74294, Lon: -157.21328,
})
if err != nil {
    panic(err)
}
fmt.Println(zone[0])
```

```bash
[ugjka@archee example]$ go run main.go
Pacific/Kiritimati
```

Uses simplified shapefile from [timezone-boundary-builder](https://github.com/evansiroky/timezone-boundary-builder/)

GeoJson Simplification done with [mapshaper](http://mapshaper.org/)

## Features

* The timezone shapefile is embedded in the build binary using go-bindata
* Supports overlapping zones
* You can load your own geojson shapefile if you want
* Sub millisecond lookup even on old hardware 

## Problems

* Shapefile is simplified using a lossy method so it may be innacurate along the borders
* This is purerly in-memory. Uses ~100MB of ram
* Nautical timezones are not included for practical reasons
