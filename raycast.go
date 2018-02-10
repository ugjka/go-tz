package gotz

import "math"

type polygon [][][]float64
type point []float64

func (p polygon) points() [][]float64 {
	return p[0]
}

func newPoint(lat, lng float64) point {
	return []float64{lng, lat}
}

func (p point) lat() float64 {
	return p[1]
}

func (p point) lng() float64 {
	return p[0]
}

// Returns whether or not the polygon is closed.
// TODO:  This can obviously be improved, but for now,
//        this should be sufficient for detecting if points
//        are contained using the raycast algorithm.
func (p polygon) IsClosed() bool {
	if len(p.points()) < 3 {
		return false
	}

	return true
}

// Returns whether or not the current Polygon contains the passed in Point.
func (p polygon) Contains(point point) bool {
	if !p.IsClosed() {
		return false
	}

	start := len(p.points()) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, p.points()[start], p.points()[end])

	for i := 1; i < len(p.points()); i++ {
		if p.intersectsWithRaycast(point, p.points()[i-1], p.points()[i]) {
			contains = !contains
		}
	}

	return contains
}

// Using the raycast algorithm, this returns whether or not the passed in point
// Intersects with the edge drawn by the passed in start and end points.
// Original implementation: http://rosettacode.org/wiki/Ray-casting_algorithm#Go
func (p polygon) intersectsWithRaycast(point point, start point, end point) bool {
	// Always ensure that the the first point
	// has a y coordinate that is less than the second point
	if start.lng() > end.lng() {

		// Switch the points if otherwise.
		start, end = end, start

	}

	// Move the point's y coordinate
	// outside of the bounds of the testing region
	// so we can start drawing a ray
	for point.lng() == start.lng() || point.lng() == end.lng() {
		newLng := math.Nextafter(point.lng(), math.Inf(1))
		point = newPoint(point.lat(), newLng)
	}

	// If we are outside of the polygon, indicate so.
	if point.lng() < start.lng() || point.lng() > end.lng() {
		return false
	}

	if start.lat() > end.lat() {
		if point.lat() > start.lat() {
			return false
		}
		if point.lat() < end.lat() {
			return true
		}

	} else {
		if point.lat() > end.lat() {
			return false
		}
		if point.lat() < start.lat() {
			return true
		}
	}

	raySlope := (point.lng() - start.lng()) / (point.lat() - start.lat())
	diagSlope := (end.lng() - start.lng()) / (end.lat() - start.lat())

	return raySlope >= diagSlope
}
