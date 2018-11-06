package gotz

type polygon []Point

func newPoint(lon, lat *float64) *Point {
	return &Point{*lon, *lat}
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
func (p polygon) contains(point *Point) bool {
	if !p.isClosed() {
		return false
	}

	start := len(p) - 1
	end := 0

	contains := p.intersectsWithRaycast(point, &p[start], &p[end])

	for i := 1; i < len(p); i++ {
		if p.intersectsWithRaycast(point, &p[i-1], &p[i]) {
			contains = !contains
		}
	}

	return contains
}

func (p polygon) intersectsWithRaycast(point, start, end *Point) bool {
	return (start.Lon > point.Lon) != (end.Lon > point.Lon) &&
		point.Lat < (end.Lat-start.Lat)*(point.Lon-start.Lon)/(end.Lon-start.Lon)+start.Lat
}
