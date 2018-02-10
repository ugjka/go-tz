# go-tz
tz-lookup (wip)

lookup timezone for a given location

```go
// Loading Zone for Line Islands, Kiritimati
	zone, err := gotz.GetZone(gotz.Point{
		Lat: 1.74294, Lng: -157.21328,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(zone)
  ```
  
  ```
  [ugjka@archee example]$ go run main.go 
Pacific/Kiritimati
```
