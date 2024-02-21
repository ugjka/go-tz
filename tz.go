// Package tz provides timezone lookup for a given location
//
// # Features
//
// * The timezone shapefile is embedded in the build binary using go-bindata
//
// * Supports overlapping zones
//
// * You can load your own geojson shapefile if you want
//
// * Sub millisecond lookup even on old hardware
//
// # Problems
//
// * The shapefile is simplified using a lossy method so it may be innacurate along the borders
//
// * This is purerly in-memory. Uses ~50MB of ram
package tz

import (
	"bytes"
	"compress/gzip"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"runtime/debug"
)

func init() {
	load()
}

//go:embed shapefile.gz
var tzShapeFile []byte

func load() {
	g, err := gzip.NewReader(bytes.NewBuffer(tzShapeFile))
	if err != nil {
		panic(err)
	}
	defer g.Close()

	if err := json.NewDecoder(g).Decode(&tzdata); err != nil {
		panic(err)
	}
	tzShapeFile = nil
	buildCenterCache()
	debug.FreeOSMemory()
}

var tzdata FeatureCollection

type centers map[string][]Point

var centerCache centers

// Point describes a location by Latitude and Longitude
type Point struct {
	Lon float64
	Lat float64
}

// ErrNoZoneFound is returned when a zone for the given point is not found in the shapefile
var ErrNoZoneFound = errors.New("no corresponding zone found in shapefile")

// ErrOutOfRange is returned when latitude exceeds 90 degrees or longitude exceeds 180 degrees
var ErrOutOfRange = errors.New("point's coordinates out of range")

// GetZone returns a slice of strings containing time zone id's for a given Point
func GetZone(p Point) (tzid []string, err error) {
	if p.Lon > 180 || p.Lon < -180 || p.Lat > 90 || p.Lat < -90 {
		return nil, ErrOutOfRange
	}
	var id string
	for _, v := range tzdata.Features {
		if v.Properties.Tzid == "" {
			continue
		}
		id = v.Properties.Tzid
		polys := v.Geometry.Coordinates
		bboxes := v.Geometry.BoundingBoxes
		holes := v.Geometry.Holes
	loop:
		for i := 0; i < len(polys); i++ {
			//Check bounding box first
			if !inBoundingBox(bboxes[i], &p) {
				continue
			}
			if polygon(polys[i]).contains(&p) {
				for j := range holes[i] {
					if polygon(holes[i][j]).contains(&p) {
						continue loop
					}
				}
				tzid = append(tzid, id)
			}
		}
	}
	if len(tzid) > 0 {
		return tzid, nil
	}
	return getClosestZone(&p)
}

func distanceFrom(p1, p2 *Point) float64 {
	d0 := (p1.Lon - p2.Lon)
	d1 := (p1.Lat - p2.Lat)
	return math.Sqrt(d0*d0 + d1*d1)
}

func getClosestZone(point *Point) (tzid []string, err error) {
	mindist := math.Inf(1)
	var winner string
	for id, v := range centerCache {
		for _, p := range v {
			tmp := distanceFrom(&p, point)
			if tmp < mindist {
				mindist = tmp
				winner = id
			}
		}
	}
	// Limit search radius
	if mindist > 2.0 {
		return getNauticalZone(point)
	}
	return append(tzid, winner), nil
}

func getNauticalZone(point *Point) (tzid []string, err error) {
	z := point.Lon / 7.5
	z = (math.Abs(z) + 1) / 2
	z = math.Floor(z)
	if z == 0 {
		return append(tzid, "Etc/GMT"), nil
	}
	if point.Lon < 0 {
		return append(tzid, fmt.Sprintf("Etc/GMT+%.f", z)), nil
	}
	return append(tzid, fmt.Sprintf("Etc/GMT-%.f", z)), nil
}

// BuildCenterCache builds centers for polygons
func buildCenterCache() {
	centerCache = make(centers)
	var tzid string
	for _, v := range tzdata.Features {
		if v.Properties.Tzid == "" {
			continue
		}
		tzid = v.Properties.Tzid
		for _, poly := range v.Geometry.Coordinates {
			centerCache[tzid] = append(centerCache[tzid], polygon(poly).centroid())
		}
	}
}

// LoadGeoJSON loads a custom GeoJSON shapefile from a Reader
func LoadGeoJSON(r io.Reader) error {
	tzdata = FeatureCollection{}
	err := json.NewDecoder(r).Decode(&tzdata)
	if err != nil {
		load()
		return err
	}
	buildCenterCache()
	debug.FreeOSMemory()
	return nil
}
