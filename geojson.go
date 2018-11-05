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
	Type        string    `json:"type"`
	Coordinates [][]Point `json:"coordinates,omitempty"`
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
		pol := make([]Point, 0)
		for _, v := range polygon.Coordinates[0] {
			pol = append(pol, Point{v[0], v[1]})
		}
		b := getBoundingBox(pol)
		g.Coordinates = append(g.Coordinates, b)
		g.Coordinates = append(g.Coordinates, pol)
		return nil
	}

	if polyType.Type == "MultiPolygon" {
		if err := json.Unmarshal(data, &multiPolygon); err != nil {
			return err
		}
		for _, poly := range multiPolygon.Coordinates {
			pol := make([]Point, 0)
			for _, v := range poly[0] {
				pol = append(pol, Point{v[0], v[1]})
			}
			b := getBoundingBox(pol)
			g.Coordinates = append(g.Coordinates, b)
			g.Coordinates = append(g.Coordinates, pol)
		}
		return nil
	}
	return nil
}
