package gotz

import (
	"fmt"
	"testing"
	"time"
)

func TestGetZone(t *testing.T) {
	//Europe/Riga
	p := Point{56.946285, 24.105078}
	start := time.Now()
	zone, err := GetZone(p)
	if err != nil {
		t.Error("Could not find Europe/Riga")
	}
	if zone[0] != "Europe/Riga" {
		t.Error("Zone not Europe/Riga but", zone[0])
	}
	fmt.Println(zone, time.Now().Sub(start))
	//Test Tokyo
	p = Point{35.6828387, 139.7594549}
	start = time.Now()
	zone, err = GetZone(p)
	if err != nil {
		t.Error("Could not find Asia/Tokyo")
	}
	if zone[0] != "Asia/Tokyo" {
		t.Error("Zone not Asia/Tokyo but", zone[0])
	}
	fmt.Println(zone, time.Now().Sub(start))
	//Tuvalu testing center cache
	p = Point{-7.768959, 178.1167698}
	start = time.Now()
	zone, err = GetZone(p)
	if err != nil {
		t.Error("Could not find Pacific/Funafuti")
	}
	if zone[0] != "Pacific/Funafuti" {
		t.Error("Zone not Pacific/Funafuti but", zone[0])
	}
	fmt.Println(zone, time.Now().Sub(start))
	//Baker Island AoE. Should error out
	p = Point{0.190165906, -176.474331436}
	start = time.Now()
	_, err = GetZone(p)
	if err == nil {
		t.Error("Baker island didn't error")
	}
	fmt.Println("Not found", time.Now().Sub(start))
}

func ExampleGetZone() {
	// Loading Zone for Line Islands, Kiritimati
	p := Point{
		Lat: 1.74294, Lon: -157.21328,
	}
	zone, err := GetZone(p)
	if err != nil {
		panic(err)
	}
	fmt.Println(zone[0])
	// Output: Pacific/Kiritimati
}
