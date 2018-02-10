// Package gotz for timezone lookup for given location
// Timezone shapefile is embedded in the build binary using go-bindata
// You can load your own geojson shapefile if you want
//
package gotz

import (
	"encoding/json"
	"errors"
	"time"
	"io"
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
	Lng float64
}

// ErrNoZoneFound when Zone for given point is not found in shapefile
var ErrNoZoneFound = errors.New("gotz: No corresponding zone found in shapefile")

// GetZone gets time.Location
func GetZone(p Point) (loc *time.Location, err error) {
	for _, v := range tzdata.Features {
		if v.Geometry.pointInZone([]float64{p.Lng, p.Lat}) {
			return time.LoadLocation(v.getZone())
		}
	}
	return nil, ErrNoZoneFound
}

// LoadGeoJSON loads custom GeoJSON shapefile
func LoadGeoJSON(r io.Reader) error {
	tzdata = FeatureCollection{}
	return json.NewDecoder(r).Decode(&tzdata)
}
