package gotz

import "math"

type polygon []Point

func newPoint(lon, lat float64) Point {
	return Point{lon, lat}
}

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

// Returns whether or not the polygon is closed.
// TODO:  This can obviously be improved, but for now,
//        this should be sufficient for detecting if points
//        are contained using the raycast algorithm.
func (p polygon) isClosed() bool {
	if len(p) < 3 {
		return false
	}

	return true
}

// Returns whether or not the current Polygon contains the passed in Point.
func (p polygon) contains(point Point) bool {
	if !p.isClosed() {
		return false
	}

	start := len(p) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, p[start], p[end])

	for i := 1; i < len(p); i++ {
		if p.intersectsWithRaycast(point, p[i-1], p[i]) {
			contains = !contains
		}
	}

	return contains
}

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func (p polygon) intersectsWithRaycast(point, start, end Point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.Lon > end.Lon {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for {
		if point.Lon != start.Lon {
			break
		}
		if point.Lon != end.Lon {
			break
		}
		newLon := math.Nextafter(point.Lon, math.Inf(1))
		point = newPoint(newLon, point.Lat)
	}

	// If we are outside of the polygon, indicate so.
	if point.Lon < start.Lon {
		return false
	}

	if point.Lon > end.Lon {
		return false
	}

	if start.Lat > end.Lat {
		if point.Lat > start.Lat {
			return false
		}
		if point.Lat < end.Lat {
			return true
		}

	} else {
		if point.Lat > end.Lat {
			return false
		}
		if point.Lat < start.Lat {
			return true
		}
	}

	raySlope := (point.Lon - start.Lon) / (point.Lat - start.Lat)
	diagSlope := (end.Lon - start.Lon) / (end.Lat - start.Lat)

	return raySlope >= diagSlope
}
