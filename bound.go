package gotz

import "math"

func getBoundingBox(points [][]float64) [][]float64 {
	if len(points) == 0 {
		return [][]float64{{0, 0}, {0, 0}}
	}

	minX := math.Inf(1)
	minY := math.Inf(1)

	maxX := math.Inf(-1)
	maxY := math.Inf(-1)

	for _, v := range points {
		minX = math.Min(minX, v[0])
		minY = math.Min(minY, v[1])

		maxX = math.Max(maxX, v[0])
		maxY = math.Max(maxY, v[1])
	}

	return [][]float64{
		{math.Min(minX, maxX), math.Min(minY, maxY)},
		{math.Max(minX, maxX), math.Max(minY, maxY)},
	}
}

func inBoundingBox(box [][]float64, point []float64) bool {

	if point[1] < box[0][1] || box[1][1] < point[1] {
		return false
	}

	if point[0] < box[0][0] || box[1][0] < point[0] {
		return false
	}

	return true
}
