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
	if zone.String() != "Europe/Riga" {
		t.Error("Zone not Europe/Riga but", zone.String())
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
