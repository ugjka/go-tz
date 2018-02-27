// Package gotz for timezone lookup for given location
// Timezone shapefile is embedded in the build binary using go-bindata
// You can load your own geojson shapefile if you want
//
package gotz

import (
	"encoding/json"
	"errors"
	"io"
	"time"
)

func init() {
	data, err := Asset("reduced/reduced.json")
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(data, &tzdata); err != nil {
		panic(err)
	}
}

var tzdata FeatureCollection

// Point is a location
type Point struct {
	Lat float64
	Lon float64
}

// ErrNoZoneFound when Zone for given point is not found in shapefile
var ErrNoZoneFound = errors.New("no corresponding zone found in shapefile")

// GetZone gets time.Location
func GetZone(p Point) (loc *time.Location, err error) {
	var tzid string
	for _, v := range tzdata.Features {
		if _, err := v.getZone(); err != nil {
			continue
		}
		for _, poly := range v.Geometry.Coordinates {
			if polygon(poly).Contains([]float64{p.Lon, p.Lat}) {
				tzid, _ = v.getZone()
				return time.LoadLocation(tzid)
			}
		}
	}
	return nil, ErrNoZoneFound
}

// LoadGeoJSON loads custom GeoJSON shapefile
func LoadGeoJSON(r io.Reader) error {
	tzdata = FeatureCollection{}
	return json.NewDecoder(r).Decode(&tzdata)
}
