package gotz

import (
	"fmt"
	"testing"
	"time"
)

func TestGetZone(t *testing.T) {
	//Europe/Riga
	p := Point{24.105078, 56.946285}
	start := time.Now()
	zone, err := GetZone(p)
	if err != nil {
		t.Error("Could not find Europe/Riga")
	}
	if len(zone) != 0 && zone[0] != "Europe/Riga" {
		t.Error("Zone not Europe/Riga but", zone[0])
	}
	fmt.Println(zone, time.Since(start))
	//Test Tokyo
	p = Point{139.7594549, 35.6828387}
	start = time.Now()
	zone, err = GetZone(p)
	if err != nil {
		t.Error("Could not find Asia/Tokyo")
	}
	if len(zone) != 0 && zone[0] != "Asia/Tokyo" {
		t.Error("Zone not Asia/Tokyo but", zone[0])
	}
	fmt.Println(zone, time.Since(start))
	//Tuvalu testing center cache
	p = Point{178.1167698, -7.768959}
	start = time.Now()
	zone, err = GetZone(p)
	if err != nil {
		t.Error("Could not find Pacific/Funafuti")
	}
	if len(zone) != 0 && zone[0] != "Pacific/Funafuti" {
		t.Error("Zone not Pacific/Funafuti but", zone[0])
	}
	fmt.Println(zone, time.Since(start))
	//Baker Island AoE. Should error out
	p = Point{-176.474331436, 0.190165906}
	start = time.Now()
	_, err = GetZone(p)
	if err == nil {
		t.Error("Baker island didn't error")
	}
	fmt.Println("Not found", time.Since(start))
}

func ExampleGetZone() {
	// Loading Zone for Line Islands, Kiritimati
	p := Point{
		Lon: -157.21328, Lat: 1.74294,
	}
	zone, err := GetZone(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(zone[0])
	// Output: Pacific/Kiritimati
}
