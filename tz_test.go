package tz

import (
	"fmt"
	"reflect"
	"testing"
)

type result struct {
	zones []string
	err   error
}

var tt = []struct {
	name  string
	point Point
	result
}{
	{
		"Riga",
		Point{24.105078, 56.946285},
		result{
			zones: []string{"Europe/Riga"},
			err:   nil,
		},
	},
	{
		"Tokyo",
		Point{139.7594549, 35.6828387},
		result{
			zones: []string{"Asia/Tokyo"},
			err:   nil,
		},
	},
	{
		"Urumqi/Shanghai",
		Point{87.319461, 43.419754},
		result{
			zones: []string{"Asia/Shanghai", "Asia/Urumqi"},
			err:   nil,
		},
	},
	{
		"Tuvalu",
		Point{178.1167698, -7.768959},
		result{
			zones: []string{"Pacific/Funafuti"},
			err:   nil,
		},
	},
	{
		"Baker Island",
		Point{-176.474331436, 0.190165906},
		result{
			zones: []string{"Etc/GMT+12"},
			err:   nil,
		},
	},
	{
		"Dubai (hole)",
		Point{56.276871, 25.276084},
		result{
			zones: []string{"Asia/Dubai"},
			err:   nil,
		},
	},
	{
		"Muscat (in hole)",
		Point{56.25347, 25.25701},
		result{
			zones: []string{"Asia/Muscat"},
			err:   nil,
		},
	},
	{
		"Kazakhstan",
		Point{76.9457275, 43.2363924},
		result{
			zones: []string{"Asia/Almaty"},
			err:   nil,
		},
	},
}

func TestGetZone(t *testing.T) {
	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			tzid, err := GetZone(tc.point)
			if err != tc.err {
				t.Errorf("expected err %v; got %v", tc.err, err)
			}
			if !reflect.DeepEqual(tzid, tc.zones) {
				t.Errorf("expected zones %v; got %v", tc.zones, tzid)
			}
		})
	}
}

func BenchmarkZones(b *testing.B) {
	b.Run("polygon centers", func(b *testing.B) {
	Loop:
		for n := 0; n < b.N; {
			for _, v := range centerCache {
				for i := range v {
					if n > b.N {
						break Loop
					}
					_, err := GetZone(v[i])
					if err != nil {
						b.Errorf("point %v did not return a zone", v[i])
					}
					n++
				}
			}
		}
	})
	b.Run("test cases", func(b *testing.B) {
	Loop:
		for n := 0; n < b.N; {
			for _, tc := range tt {
				if n > b.N {
					break Loop
				}
				GetZone(tc.point)
				n++
			}

		}
	})
}

func TestNautical(t *testing.T) {
	tt := []struct {
		lon  float64
		zone string
	}{
		{-180, "Etc/GMT+12"},
		{180, "Etc/GMT-12"},
		{-172.5, "Etc/GMT+12"},
		{172.5, "Etc/GMT-12"},
		{-172, "Etc/GMT+11"},
		{172, "Etc/GMT-11"},
		{0, "Etc/GMT"},
		{7.49, "Etc/GMT"},
		{-7.49, "Etc/GMT"},
		{7.5, "Etc/GMT-1"},
		{-7.5, "Etc/GMT+1"},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%f %s", tc.lon, tc.zone), func(t *testing.T) {
			z, _ := getNauticalZone(&Point{Lat: 0, Lon: tc.lon})
			if z[0] != tc.zone {
				t.Errorf("expected %s got %s", tc.zone, z[0])
			}
		})
	}
}

func TestOutOfRange(t *testing.T) {
	tt := []struct {
		p   Point
		err error
	}{
		{Point{180, 0}, nil},
		{Point{-180, 0}, nil},
		{Point{0, 90}, nil},
		{Point{0, -90}, nil},
		{Point{181, 0}, ErrOutOfRange},
		{Point{-181, 0}, ErrOutOfRange},
		{Point{0, 91}, ErrOutOfRange},
		{Point{0, -91}, ErrOutOfRange},
	}
	for _, tc := range tt {
		t.Run(fmt.Sprintf("%f %f", tc.p.Lon, tc.p.Lat), func(t *testing.T) {
			_, err := GetZone(tc.p)
			if err != tc.err {
				t.Errorf("expected error %v got %v", tc.err, err)
			}
		})
	}
}
