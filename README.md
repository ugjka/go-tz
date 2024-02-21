# go-tz

Time zone lookup by Lon and Lat

## Usage

```go
import "github.com/ugjka/go-tz/v2"
```

### Example

```go
// Loading Zone for Line Islands, Kiritimati
zone, err := tz.GetZone(tz.Point{
    Lon: -157.21328, Lat: 1.74294,
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

- The timezone shapefile is embedded in the build binary
- Supports overlapping zones
- You can load your custom geojson shapefile if you want
- Sub-millisecond lookup even on old hardware

## Problems

- Shapefile is simplified using a lossy method so it may be inaccurate along the borders
- This is purely in-memory. Uses ~50MB of ram

## Licenses

The code used to look up the timezone for a location is licensed under the [MIT License](https://opensource.org/licenses/MIT).

The data in the timezone shapefile is licensed under the [Open Data Commons Open Database License (ODbL)](http://opendatacommons.org/licenses/odbl/).
