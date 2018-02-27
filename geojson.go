package gotz

import (
	"encoding/json"
	"errors"
)

//ErrNoTZID when no TZID found for region
var ErrNoTZID = errors.New("tzid for feature not found")

// FeatureCollection ...
type FeatureCollection struct {
	featureCollection
}

type featureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

// Feature ...
type Feature struct {
	feature
}

type feature struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Geometry   Geometry          `json:"geometry"`
	Properties map[string]string `json:"properties"`
}

// Geometry ...
type Geometry struct {
	geometry
}

type geometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates,omitempty"`
}

func (f *Feature) getZone() (string, error) {
	if v, ok := f.Properties["TZID"]; ok {
		return v, nil
	}
	if v, ok := f.Properties["tzid"]; ok {
		return v, nil
	}
	return "", ErrNoTZID
}

// UnmarshalJSON for polygons and multipolygons
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
		g.Coordinates = append(g.Coordinates, polygon.Coordinates)
		return nil
	}

	if polyType.Type == "MultiPolygon" {
		if err := json.Unmarshal(data, &multiPolygon); err != nil {
			return err
		}
		g.Coordinates = multiPolygon.Coordinates
		return nil
	}
	return nil
}
