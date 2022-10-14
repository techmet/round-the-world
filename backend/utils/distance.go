package utils

import (
	"math"

	"github.com/techmet/round-the-world/models"
)

func GetDistanceFromLatLonInKm(city1 models.City, city2 models.City) float64 {
	var R float64 = 6371
	dLat := deg2rad(city2.Location.Lat - city1.Location.Lat)
	dLon := deg2rad(city2.Location.Lon - city1.Location.Lon)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(deg2rad(city1.Location.Lat))*math.Cos(deg2rad(city2.Location.Lat))*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func deg2rad(deg float64) float64 {
	return (deg / 180) * math.Pi
}
