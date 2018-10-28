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
func GetZone(p Point) (tzid []string, err error) {
	var id string
	for _, v := range tzdata.Features {
		if id, err = v.getTZID(); err != nil {
			continue
		}
		polys := v.Geometry.Coordinates
		for i := 0; i < len(polys); i += 2 {
			//Check bounding box first
			//Massive speedup
			if !inBoundingBox(polys[i][0], []float64{p.Lon, p.Lat}) {
				continue
			}
			if polygon(polys[i+1]).contains([]float64{p.Lon, p.Lat}) {
				tzid = append(tzid, id)
			}
		}
	}
	if len(tzid) > 0 {
		return tzid, nil
	}
	return getClosestZone(p)
}

func distanceFrom(p1, p2 []float64) float64 {
	d0 := (p1[0] - p2[0])
	d1 := (p1[1] - p2[1])
	return math.Sqrt(d0*d0 + d1*d1)
}

func getClosestZone(point Point) (tzid []string, err error) {
	mindist := math.Inf(1)
	var winner string
	for id, v := range centerCache {
		for _, p := range v {
			tmp := distanceFrom(p, []float64{point.Lon, point.Lat})
			if tmp < mindist {
				mindist = tmp
				winner = id
			}
		}
	}
	// Limit search radius
	if mindist > 2.0 {
		return tzid, ErrNoZoneFound
	}
	return append(tzid, winner), nil
}

//BuildCenterCache builds centers for polygons
func buildCenterCache() {
	centerCache = make(centers)
	var tzid string
	var err error
	for _, v := range tzdata.Features {
		if tzid, err = v.getTZID(); err != nil {
			continue
		}
		for i, poly := range v.Geometry.Coordinates {
			// ignore bounding boxes
			if i%2 == 0 {
				continue
			}
			centerCache[tzid] = append(centerCache[tzid], polygon(poly).centroid())
		}
	}
}

// LoadGeoJSON loads custom GeoJSON shapefile
func LoadGeoJSON(r io.Reader) error {
	tzdata = FeatureCollection{}
	err := json.NewDecoder(r).Decode(&tzdata)
	if err != nil {
		return err
	}
	buildCenterCache()
	return nil
}
