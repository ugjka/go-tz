// Program for simplifying the GeoJSON polygons
// Not great but it will suffice for now
// Pull requests are welcome
// Needs smarter algo
package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	"github.com/yrsh/simplify-go"
)

const (
	// IGNORE polygons with points less than this
	IGNORE = 100
)

func main() {
	var data FeatureCollection
	in, err := os.Open("./dist/combined.json")
	if err != nil {
		panic(err)
	}
	err = json.NewDecoder(in).Decode(&data)
	if err != nil {
		panic(err)
	}
	os.Mkdir("./reduced", 0755)
	out, err := os.OpenFile("./reduced/reduced.json", os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0655)
	if err != nil {
		panic(err)
	}
	err = json.NewEncoder(out).Encode(data)
	if err != nil {
		panic(err)
	}
	in.Close()
	out.Close()
}

// FeatureCollection ...
type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

// Feature ...
type Feature struct {
	Type       string            `json:"type"`
	ID         string            `json:"id"`
	Geometry   Geometry          `json:"geometry"`
	Properties map[string]string `json:"properties"`
}

// Geometry ...
type Geometry struct {
	Type        string          `json:"type"`
	Coordinates [][][][]float64 `json:"coordinates,omitempty"`
	Geometries  []*Geometry     `json:"geometries,omitempty"`
}

// UnmarshalJSON where simplification happens
func (g *Geometry) UnmarshalJSON(data []byte) (err error) {
	var polyType struct {
		Type       string      `json:"type"`
		Geometries []*Geometry `json:"geometries,omitempty"`
	}
	if err := json.Unmarshal(data, &polyType); err != nil {
		return err
	}
	g.Geometries = polyType.Geometries
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
		for i, v := range g.Coordinates[0] {
			g.Coordinates[0][i] = simplify(v)
		}
		return nil
	}

	if polyType.Type == "MultiPolygon" {
		if err := json.Unmarshal(data, &multiPolygon); err != nil {
			return err
		}
		g.Coordinates = multiPolygon.Coordinates
		for i, v := range g.Coordinates[0] {
			g.Coordinates[0][i] = simplify(v)
		}
		return nil
	}
	return nil
}

// PolyLine simplification using
// https://github.com/yrsh/simplify-go
func simplify(v [][]float64) (out [][]float64) {
	if len(v) < IGNORE {
		return v
	}
	for i := 0.0005; ; i = i + 0.0005 {
		slice := simplifier.Simplify(v, i, true)
		if len(slice) < int(math.Sqrt(float64(len(v)*2))*3)+IGNORE {
			if len(slice) < 20 {
				fmt.Printf("%d, ", len(slice))
			}
			return slice
		}
	}
}
