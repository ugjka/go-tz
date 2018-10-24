// Package gotz for timezone lookup for given location
// Timezone shapefile is embedded in the build binary using go-bindata
// You can load your own geojson shapefile if you want
//
package gotz

import (
	"encoding/json"
	"errors"
	"io"
	"math"
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
	buildCenterCache()
}

var tzdata FeatureCollection

type centers map[string][][]float64

var centerCache centers

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
		if tzid, err = v.getZone(); err != nil {
			continue
		}
		for _, poly := range v.Geometry.Coordinates {
			if polygon(poly).Contains([]float64{p.Lon, p.Lat}) {
				return time.LoadLocation(tzid)
			}
		}
	}
	return getClosestZone(p)
}

func distanceFrom(p1, p2 []float64) float64 {
	d0 := (p1[0] - p2[0])
	d1 := (p1[1] - p2[1])
	return math.Sqrt(d0*d0 + d1*d1)
}

func getClosestZone(point Point) (loc *time.Location, err error) {
	mindist := math.MaxFloat64
	var winner string
	for tzid, v := range centerCache {
		for _, p := range v {
			tmp := distanceFrom(p, []float64{point.Lon, point.Lat})
			tmp = math.Abs(tmp)
			if tmp < mindist {
				mindist = tmp
				winner = tzid
			}
		}
	}
	// Limit search radius
	if mindist > 2.0 {
		return nil, ErrNoZoneFound
	}
	return time.LoadLocation(winner)
}

//BuildCenterCache builds centers for polygons
func buildCenterCache() {
	centerCache = make(centers)
	var tzid string
	var err error
	for _, v := range tzdata.Features {
		if tzid, err = v.getZone(); err != nil {
			continue
		}
		for _, poly := range v.Geometry.Coordinates {
			if _, ok := centerCache[tzid]; !ok {
				centerCache[tzid] = make([][]float64, 0)
			}
			centerCache[tzid] = append(centerCache[tzid], polygon(poly).Centroid())
		}
	}
}

// LoadGeoJSON loads custom GeoJSON shapefile
func LoadGeoJSON(r io.Reader) error {
	tzdata = FeatureCollection{}
	return json.NewDecoder(r).Decode(&tzdata)
}
