package models

type LocationDetails struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}
type City struct {
	ID                 string          `json:"id"`
	Name               string          `json:"name"`
	Location           LocationDetails `json:"location"`
	CountryName        string          `json:"countryName"`
	Iata               string          `json:"iata"`
	Rank               int             `json:"rank"`
	CountryID          string          `json:"countryId"`
	Dest               interface{}     `json:"dest"`
	Airports           []string        `json:"airports"`
	Images             []string        `json:"images"`
	Popularity         float64         `json:"popularity"`
	RegID              string          `json:"regId"`
	ContID             string          `json:"contId"`
	Con                int             `json:"con"`
	NeighbouringCities NeighbouringContinentCityMap
}

type SearchCity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type CityMap map[string]*City
type ContinentCityMap map[string][]*City
type NeighbouringContinentCityMap map[string][]*NeighbouringCity

type NeighbouringCity struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	CountryName string  `json:"countryName"`
	Distance    float64 `json:"distance"`
	Continent   string  `json:"continent"`
}

type TripDetails struct {
	ID            string            `json:"id"`
	Path          []string          `json:"path"`
	Coordinates   []LocationDetails `json:"locationDetails"`
	TotalDistance float64           `json:"totalDistance"`
}
