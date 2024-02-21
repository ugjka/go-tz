package tz

import (
	"math"
)

type polygon []Point

// func newPoint(lon, lat *float64) *Point {
// 	return &Point{*lon, *lat}
// }

func (p polygon) centroid() Point {
	x := 0.0
	y := 0.0
	numPoints := float64(len(p))
	for _, p := range p {
		x += p.Lon
		y += p.Lat
	}
	return Point{x / numPoints, y / numPoints}
}

func (p polygon) isClosed() bool {
	return len(p) > 3
}

// Returns whether or not the current Polygon contains the passed in Point.
func (p polygon) contains(point *Point) bool {
	if !p.isClosed() {
		return false
	}

	start := len(p) - 1
	end := 0

	contains := intersectsWithRaycast(point, &p[start], &p[end])

	for i := 1; i < len(p); i++ {
		if intersectsWithRaycast(point, &p[i-1], &p[i]) {
			contains = !contains
		}
	}

	return contains
}

// https://rosettacode.org/wiki/Ray-casting_algorithm#Go
func intersectsWithRaycast(point, start, end *Point) bool {
	if start.Lat > end.Lat {
		start, end = end, start
	}
	for point.Lat == start.Lat || point.Lat == end.Lat {
		point.Lat = math.Nextafter(point.Lat, math.Inf(1))
	}
	if point.Lat < start.Lat || point.Lat > end.Lat {
		return false
	}
	if start.Lon > end.Lon {
		if point.Lon > start.Lon {
			return false
		}
		if point.Lon < end.Lon {
			return true
		}
	} else {
		if point.Lon > end.Lon {
			return false
		}
		if point.Lon < start.Lon {
			return true
		}
	}
	return (point.Lat-start.Lat)/(point.Lon-start.Lon) >= (end.Lat-start.Lat)/(end.Lon-start.Lon)
}
