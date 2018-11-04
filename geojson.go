package gotz

import (
	"encoding/json"
	"errors"
)

var errNoTZID = errors.New("tzid for feature not found")

type FeatureCollection struct {
	featureCollection
}

type featureCollection struct {
	Features []*Feature `json:"features"`
}

type Feature struct {
	feature
}

type feature struct {
	Geometry   Geometry `json:"geometry"`
	Properties struct {
		Tzid string `json:"tzid"`
	} `json:"properties"`
}

type Geometry struct {
	geometry
}

type geometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates,omitempty"`
}

func (g *Geometry) UnmarshalJSON(data []byte) (err error) {
	var polyType struct {
		Type       string      `json:"type"`
		Geometries []*Geometry `json:"geometries,omitempty"`
	}
	if err := json.Unmarshal(data, &polyType); err != nil {
		return err
	}
	g.Type = "MultiPolygon"
	var polygon struct {
		Coordinates [][][]float64 `json:"coordinates,omitempty"`
	}

	var multiPolygon struct {
		Coordinates [][][][]float64 `json:"coordinates,omitempty"`
	}

	if polyType.Type == "Polygon" {
		if err := json.Unmarshal(data, &polygon); err != nil {
			return err
		}
		//Create a bounding box
		b := make([][][]float64, 1)
		b[0] = getBoundingBox(polygon.Coordinates[0])
		g.Coordinates = append(g.Coordinates, b)
		g.Coordinates = append(g.Coordinates, polygon.Coordinates)
		return nil
	}

	if polyType.Type == "MultiPolygon" {
		if err := json.Unmarshal(data, &multiPolygon); err != nil {
			return err
		}
		for _, poly := range multiPolygon.Coordinates {
			//Create a bounding box
			b := make([][][]float64, 1)
			b[0] = getBoundingBox(poly[0])
			g.Coordinates = append(g.Coordinates, b)
			g.Coordinates = append(g.Coordinates, poly)
		}
		return nil
	}
	return nil
}
