package tz

import "math"

func getBoundingBox(points []Point) []Point {
	if len(points) == 0 {
		return []Point{{0, 0}, {0, 0}}
	}

	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for _, v := range points {
		minX = math.Min(minX, v.Lon)
		minY = math.Min(minY, v.Lat)

		maxX = math.Max(maxX, v.Lon)
		maxY = math.Max(maxY, v.Lat)
	}

	return []Point{
		{math.Min(minX, maxX), math.Min(minY, maxY)},
		{math.Max(minX, maxX), math.Max(minY, maxY)},
	}
}

func inBoundingBox(box []Point, point *Point) bool {

	if point.Lat < box[0].Lat || box[1].Lat < point.Lat {
		return false
	}

	if point.Lon < box[0].Lon || box[1].Lon < point.Lon {
		return false
	}

	return true
}
