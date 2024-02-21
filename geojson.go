package tz

import (
	"encoding/json"
)

// FeatureCollection ...
type FeatureCollection struct {
	featureCollection
}

type featureCollection struct {
	Features []*Feature
}

// Feature ...
type Feature struct {
	feature
}

type feature struct {
	Geometry   Geometry
	Properties struct {
		Tzid string
	}
}

// Geometry ...
type Geometry struct {
	geometry
}

type geometry struct {
	Coordinates   [][]Point
	Holes         [][][]Point
	BoundingBoxes [][]Point
}

var jPolyType struct {
	Type       string
	Geometries []*Geometry
}

var jPolygon struct {
	Coordinates [][][]float64
}

var jMultiPolygon struct {
	Coordinates [][][][]float64
}

// UnmarshalJSON ...
func (g *Geometry) UnmarshalJSON(data []byte) (err error) {
	if err := json.Unmarshal(data, &jPolyType); err != nil {
		return err
	}

	if jPolyType.Type == "Polygon" {
		if err := json.Unmarshal(data, &jPolygon); err != nil {
			return err
		}
		pol := make([]Point, len(jPolygon.Coordinates[0]))
		for i, v := range jPolygon.Coordinates[0] {
			pol[i].Lon = v[0]
			pol[i].Lat = v[1]
		}
		b := getBoundingBox(pol)
		g.BoundingBoxes = append(g.BoundingBoxes, b)
		g.Coordinates = append(g.Coordinates, pol)

		var holes [][]Point
		if len(jPolygon.Coordinates) > 1 {
			holes = make([][]Point, len(jPolygon.Coordinates)-1)
			for i := 1; i < len(jPolygon.Coordinates); i++ {
				hole := make([]Point, len(jPolygon.Coordinates[i]))
				for ii, v := range jPolygon.Coordinates[i] {
					hole[ii].Lon = v[0]
					hole[ii].Lat = v[1]
				}
				holes[i-1] = hole
			}
		}
		g.Holes = append(g.Holes, holes)
		return nil
	}

	if jPolyType.Type == "MultiPolygon" {
		if err := json.Unmarshal(data, &jMultiPolygon); err != nil {
			return err
		}
		g.BoundingBoxes = make([][]Point, len(jMultiPolygon.Coordinates))
		g.Coordinates = make([][]Point, len(jMultiPolygon.Coordinates))
		for j, poly := range jMultiPolygon.Coordinates {
			pol := make([]Point, len(poly[0]))
			for i, v := range poly[0] {
				pol[i].Lon = v[0]
				pol[i].Lat = v[1]
			}
			b := getBoundingBox(pol)
			g.BoundingBoxes[j] = b
			g.Coordinates[j] = pol

			var holes [][]Point
			if len(poly) > 1 {
				holes = make([][]Point, len(poly)-1)
				for i := 1; i < len(poly); i++ {
					hole := make([]Point, len(poly[i]))
					for ii, v := range poly[i] {
						hole[ii].Lon = v[0]
						hole[ii].Lat = v[1]
					}
					holes[i-1] = hole
				}
			}
			g.Holes = append(g.Holes, holes)
		}
		return nil
	}
	return nil
}
